package reader

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/password"
	"github.com/aaronland/go-feed-reader/tables"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/aaronland/go-sql-pagination"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"io"
	"log"
	"strings"
	"sync"
)

type FeedReader struct {
	login.Provider // which implements user.UserDB
	database       *database.SQLiteDatabase
	feeds          sqlite.Table
	items          sqlite.Table
	search         sqlite.Table
	users          sqlite.Table
	user_feeds     sqlite.Table
	user_items     sqlite.Table
	cfg            login.Config
	mu             *sync.Mutex
}

type FeedsResponse struct {
	Feeds      []*gofeed.Feed
	Pagination pagination.Pagination
}

type ItemsResponse struct {
	Items      []*gofeed.Item
	Pagination pagination.Pagination
}

type ListItemsOptions struct {
	FeedURL  string
	IsRead   bool
	IsUnread bool
}

func NewDefaultListItemsOptions() *ListItemsOptions {

	opts := ListItemsOptions{
		FeedURL:  "",
		IsRead:   false,
		IsUnread: false,
	}

	return &opts
}

func NewFeedReader(dsn string) (*FeedReader, error) {

	db, err := database.NewDBWithDriver("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	err = db.LiveHardDieFast()

	if err != nil {
		return nil, err
	}

	f, err := tables.NewFeedsTableWithDatabase(db)

	if err != nil {
		return nil, err
	}

	i, err := tables.NewItemsTableWithDatabase(db)

	if err != nil {
		return nil, err
	}

	s, err := tables.NewSearchTableWithDatabase(db)

	if err != nil {
		return nil, err
	}

	u, err := tables.NewUsersTableWithDatabase(db)

	if err != nil {
		return nil, err
	}

	uf, err := tables.NewUserFeedsTableWithDatabase(db)

	if err != nil {
		return nil, err
	}

	ui, err := tables.NewUserItemsTableWithDatabase(db)

	if err != nil {
		return nil, err
	}

	ck_cfg, err := NewFRCookieConfig()

	if err != nil {
		return nil, err
	}

	url_cfg, err := login.NewDefaultURLConfig()

	if err != nil {
		return nil, err
	}

	cfg, err := NewFRConfig(ck_cfg, url_cfg)

	if err != nil {
		return nil, err
	}

	mu := new(sync.Mutex)

	fr := FeedReader{
		database:   db,
		feeds:      f,
		items:      i,
		search:     s,
		users:      u,
		user_feeds: uf,
		user_items: ui,
		mu:         mu,
		cfg:        cfg,
	}

	return &fr, nil
}

type FRConfig struct {
	login.Config
	cookie login.CookieConfig
	url    login.URLConfig
}

func (c *FRConfig) Cookie() login.CookieConfig {
	return c.cookie
}

func (c *FRConfig) URL() login.URLConfig {
	return c.url
}

type FRCookieConfig struct {
	login.CookieConfig
	salt   string
	secret string
	name   string
}

func (c *FRCookieConfig) Salt() string {
	return c.salt
}

func (c *FRCookieConfig) Secret() string {
	return c.secret
}

func (c *FRCookieConfig) Name() string {
	return c.name
}

func NewFRConfig(ck_cfg login.CookieConfig, url_cfg login.URLConfig) (login.Config, error) {

	cfg := FRConfig{
		cookie: ck_cfg,
		url:    url_cfg,
	}

	return &cfg, nil
}

func NewFRCookieConfig() (login.CookieConfig, error) {

	cfg := FRCookieConfig{
		salt:   "salty",  // PLEASE FIX ME
		secret: "cookie", // PLEASE FIX ME
		name:   "fr",     // PLEASE FIX ME
	}

	return &cfg, nil
}

// login.Provider methods

func (fr *FeedReader) Config() login.Config {
	return fr.cfg
}

// user.User methods

func (fr *FeedReader) GetUserById(id string) (user.User, error) {

	return fr.getUser("id", id)
}

func (fr *FeedReader) GetUserByEmail(email string) (user.User, error) {

	return fr.getUser("email", email)
}

func (fr *FeedReader) GetUserByUsername(name string) (user.User, error) {

	return fr.getUser("name", name)
}

func (fr *FeedReader) getUser(col string, ref string) (user.User, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT id, name, email, password FROM %s WHERE `%s` = ?", fr.users.Name(), col)
	row := conn.QueryRow(query, ref)

	var id string
	var username string
	var email string
	var digest string

	err = row.Scan(&id, &username, &email, &digest)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, &user.ErrNoUser{}
		}

		return nil, err
	}

	salt := "FIXME"
	pswd, err := password.NewBCryptPasswordFromDigest(digest, salt)

	if err != nil {
		return nil, err
	}

	u, err := user.NewDefaultUserWithID(id, username, email, pswd)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (fr *FeedReader) AddUser(u user.User) error {

	return fr.users.IndexRecord(fr.database, u)
}

// feed reader methods

func (fr *FeedReader) AddFeedForUser(u user.User, feed_url string) (*gofeed.Feed, error) {

	feed, err := fr.ParseFeedURL(feed_url)

	if err != nil {
		return nil, err
	}

	err = fr.IndexFeedForUser(u, feed)

	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (fr *FeedReader) DumpFeedsForUser(u user.User, wr io.Writer) error {

	cb := func(f *gofeed.Feed) error {
		wr.Write([]byte(f.FeedLink + "\n"))
		return nil
	}

	err := fr.ListFeedsAllForUser(u, cb)

	if err != nil {
		log.Fatal()
	}
	return nil
}

func (fr *FeedReader) GetFeedByLinkForUser(u user.User, link string) (*gofeed.Feed, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`SELECT f.body FROM %s f, %s uf
		WHERE f.link = uf.feed_link
		AND uf.feed_link = ?
		AND uf.user_id = ?`, fr.feeds.Name(), fr.user_feeds.Name())

	row := conn.QueryRow(q, link, u.Id())
	return DatabaseRowToFeed(row)
}

func (fr *FeedReader) GetFeedByItemGUIDForUser(u user.User, guid string) (*gofeed.Feed, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`SELECT f.body FROM %s f, %s i
	  	WHERE f.link = i.feed_link
		AND i.item_guid = ?
		AND i.user_id = ?`,
		fr.feeds.Name(), fr.user_items.Name())

	row := conn.QueryRow(q, guid, u.Id())
	return DatabaseRowToFeed(row)
}

func (fr *FeedReader) GetItemByGUIDForUser(u user.User, guid string) (*gofeed.Item, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`SELECT body FROM %s i, %s ui
	  	WHERE i.guid = ui.guid
		AND ui.guid = ?
		AND ui.user_id = ?`,
		fr.items.Name(), fr.user_items.Name())

	row := conn.QueryRow(q, guid, u.Id())
	return DatabaseRowToFeedItem(row)
}

func (fr *FeedReader) SearchForUser(u user.User, q string, opts pagination.PaginatedOptions) (*ItemsResponse, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	// https://www.sqlite.org/fts5.html

	query := fmt.Sprintf(`SELECT s.feed, s.guid FROM %s s, %s uf
	      WHERE uf.feed_link = s.feed
	      AND uf.user_id = ?
	      AND %s MATCH ?
	      ORDER BY rank`, fr.search.Name(), fr.user_feeds.Name(), fr.search.Name())

	rsp, err := pagination.QueryPaginated(conn, opts, query, u.Id(), q)

	if err != nil {
		return nil, err
	}

	guids := make([][]string, 0)

	rows := rsp.Rows()
	defer rows.Close()

	pg := rsp.Pagination()

	for rows.Next() {

		var feed string
		var guid string

		err = rows.Scan(&feed, &guid)

		if err != nil {
			return nil, err
		}

		guids = append(guids, []string{feed, guid})
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	// please do this concurrently

	items := make([]*gofeed.Item, 0)

	for _, g := range guids {

		feed := g[0]
		guid := g[1]

		query := fmt.Sprintf(`SELECT i.body FROM %s i, %s ui
			WHERE i.guid == ui.item_guid
			AND ui.user_id = ?
			AND ui.feed = ? AND ui.guid = ?`, fr.items.Name(), fr.user_items.Name())

		row := conn.QueryRow(query, u.Id(), feed, guid)
		item, err := DatabaseRowToFeedItem(row)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	r := ItemsResponse{
		Items:      items,
		Pagination: pg,
	}

	return &r, nil
}

func (fr *FeedReader) RemoveFeedForUser(u user.User, f *gofeed.Feed) error {

	conn, err := fr.database.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	sql_items := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND feed_link = ?", fr.user_items.Name())
	sql_feeds := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND link  = ?", fr.user_feeds.Name())

	queries := []string{
		sql_items,
		sql_feeds,
	}

	for _, q := range queries {

		_, err := conn.Exec(q, u.Id(), f.Link)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (fr *FeedReader) PruneFeed(f *gofeed.Feed) error {

	conn, err := fr.database.Conn()

	if err != nil {
		return err
	}

	sql_count := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE feed_link = ?", fr.user_feeds.Name())
	row, err := conn.Query(sql_count, f.Link)

	if err != nil {
		return err
	}

	var count int32
	err = row.Scan(&count)

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	if count == 0 {

		sql_items := fmt.Sprintf("DELETE FROM %s WHERE feed_link = ?", fr.items.Name())
		sql_feeds := fmt.Sprintf("DELETE FROM %s WHERE link = ?", fr.feeds.Name())
		sql_search := fmt.Sprintf("DELETE FROM %s WHERE feed = ?", fr.search.Name())

		queries := []string{
			sql_feeds,
			sql_items,
			sql_search,
		}

		for _, q := range queries {

			_, err := conn.Exec(q, f.Link)

			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	return tx.Commit()
}

func (fr *FeedReader) ListItemsForUser(u user.User, ls_opts *ListItemsOptions, pg_opts pagination.PaginatedOptions) (*ItemsResponse, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	conditions = append(conditions, "ui.user_id=?")
	args = append(args, u.Id())

	if ls_opts.FeedURL != "" {

		conditions = append(conditions, "i.feed = ?")
		args = append(args, ls_opts.FeedURL)
	}

	extra := ""

	if len(conditions) > 0 {
		extra = fmt.Sprintf("AND %s", strings.Join(conditions, " AND "))
	}

	q := fmt.Sprintf(`SELECT i.body FROM %s i, %s ui
		WHERE i.guid = ui.item_guid
	  	%s
	  	ORDER BY i.published ASC, i.updated ASC`, fr.items.Name(), fr.user_items.Name(), extra)

	rsp, err := pagination.QueryPaginated(conn, pg_opts, q, args...)

	if err != nil {
		return nil, err
	}

	rows := rsp.Rows()
	defer rows.Close()

	items, err := DatabaseRowsToFeedItems(rows)

	if err != nil {
		return nil, err
	}

	r := ItemsResponse{
		Items:      items,
		Pagination: rsp.Pagination(),
	}

	return &r, nil
}

func (fr *FeedReader) ListFeedsAllForUser(u user.User, feed_cb func(f *gofeed.Feed) error) error {

	query := fmt.Sprintf(`SELECT f.* FROM %s f, %s uf
		WHERE f.link = uf.feed_link
		AND uf.user_id = ?`, fr.feeds.Name(), fr.user_feeds.Name())

	opts := pagination.NewDefaultPaginatedOptions()
	opts.PerPage(100)

	return fr.listFeedsAll(opts, feed_cb, query, u.Id())
}

func (fr *FeedReader) ListFeedsAll(feed_cb func(f *gofeed.Feed) error) error {

	query := fmt.Sprintf("SELECT body FROM %s", fr.feeds.Name())
	log.Println("QUERY", query)

	opts := pagination.NewDefaultPaginatedOptions()
	opts.PerPage(100)

	return fr.listFeedsAll(opts, feed_cb, query)
}

func (fr *FeedReader) listFeedsAll(opts pagination.PaginatedOptions, feed_cb func(f *gofeed.Feed) error, query string, args ...interface{}) error {

	log.Println("LIST FEEDS ALL...", query)

	cb := func(r pagination.PaginatedResponse) error {

		rows := r.Rows()
		defer rows.Close()

		feeds, err := DatabaseRowsToFeeds(rows)

		if err != nil {
			return err
		}

		for _, feed := range feeds {

			err := feed_cb(feed)

			if err != nil {
				return err
			}
		}

		return nil
	}

	conn, err := fr.database.Conn()

	if err != nil {
		return nil
	}

	return pagination.QueryPaginatedAll(conn, opts, cb, query, args...)
}

func (fr *FeedReader) ListFeedsForUser(u user.User, pg_opts pagination.PaginatedOptions) (*FeedsResponse, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`SELECT f.body FROM %s f, %s u
	  	WHERE u.user_id = ? AND u.feed_link = f.link
		ORDER BY f.updated DESC`, fr.feeds.Name(), fr.user_feeds.Name())

	rsp, err := pagination.QueryPaginated(conn, pg_opts, q, u.Id())

	if err != nil {
		return nil, err
	}

	rows := rsp.Rows()
	defer rows.Close()

	pg := rsp.Pagination()

	feeds := make([]*gofeed.Feed, 0)

	for rows.Next() {

		var body string
		err = rows.Scan(&body)

		if err != nil {
			return nil, err
		}

		var f gofeed.Feed

		err := json.Unmarshal([]byte(body), &f)

		if err != nil {
			return nil, err
		}

		feeds = append(feeds, &f)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	r := FeedsResponse{
		Feeds:      feeds,
		Pagination: pg,
	}

	return &r, nil
}

func (fr *FeedReader) RefreshFeeds() error {

	log.Println("REFRESHING...")

	fr.mu.Lock()

	defer func() {
		fr.mu.Unlock()
	}()

	// check last update here...

	cb := func(feed *gofeed.Feed) error {

		log.Println("REFRESH FEED")

		err := fr.RefreshFeed(feed)

		if err != nil {
			log.Println("REFRESH FEED ERROR", err)
			return err
		}

		return nil
	}

	return fr.ListFeedsAll(cb)
}

func (fr *FeedReader) RefreshFeed(feed *gofeed.Feed) error {

	f2, err := fr.ParseFeedURL(feed.FeedLink)

	if err != nil {
		return err
	}

	err = fr.IndexFeed(f2)

	if err != nil {
		return err
	}

	log.Println("REFRESH FOR USERS")
	err = fr.RefreshFeedForUsers(f2)

	if err != nil {
		log.Println("FUUUUUU")
		return err
	}

	return nil
}

func (fr *FeedReader) RefreshFeedForUsers(f *gofeed.Feed) error {

	log.Println("REFRESH FOR USERS")

	conn, err := fr.database.Conn()

	if err != nil {
		return err
	}

	cb := func(r pagination.PaginatedResponse) error {

		log.Println("REFRESH FOR USERS PAGINATED RESPONSE")

		rows := r.Rows()
		defer rows.Close()

		for rows.Next() {

			log.Println("NEXT")
			var user_id string
			err := rows.Scan(&user_id)

			if err != nil {
				log.Println("NEXT ERROR", err)
				return err
			}

			u, err := fr.GetUserById(user_id)

			if err != nil {
				log.Println("NEXT USER ERROR", err, user_id)
				return err
			}

			log.Println("INDEX FEED FOR USER")
			err = fr.IndexFeedForUser(u, f)

			if err != nil {
				return err
			}
		}

		log.Println("NEXT DONE")
		return nil
	}

	query := fmt.Sprintf("SELECT user_id FROM %s WHERE feed_link=?", fr.user_feeds.Name())
	log.Println("QUERY", query)

	opts := pagination.NewDefaultPaginatedOptions()
	opts.PerPage(100)

	return pagination.QueryPaginatedAll(conn, opts, cb, query, f.Link)
}

func (fr *FeedReader) ParseFeedURL(feed_url string) (*gofeed.Feed, error) {

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feed_url)

	if err != nil {
		return nil, err
	}

	feed.FeedLink = feed_url // this shouldn't be necessary but... you know, is (20180407/thisisaaronland)
	return feed, nil
}

func (fr *FeedReader) IndexFeedForUser(u user.User, feed *gofeed.Feed) error {

	err := fr.IndexFeed(feed)

	if err != nil {
		return err
	}

	uf := tables.UserFeed{
		Feed: feed,
		User: u,
	}

	err = fr.user_feeds.IndexRecord(fr.database, &uf)

	if err != nil {
		return err
	}

	for _, item := range feed.Items {

		ui := tables.UserItem{
			User: u,
			Feed: feed,
			Item: item,
		}

		err := fr.user_items.IndexRecord(fr.database, &ui)

		if err != nil {
			return err
		}

		// something something search here...
	}

	return nil
}

func (fr *FeedReader) IndexFeed(feed *gofeed.Feed) error {

	items := feed.Items
	feed.Items = nil

	rec := tables.FeedRecord{
		Feed: feed,
	}

	err := fr.feeds.IndexRecord(fr.database, &rec)

	if err != nil {
		return err
	}

	for _, item := range items {

		rec := tables.ItemRecord{
			Feed: feed,
			Item: item,
		}

		err = fr.items.IndexRecord(fr.database, &rec)

		if err != nil {
			return err
		}

	}

	return nil
}
