package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronland/go-feed-reader/tables"
	"github.com/aaronland/go-sql-pagination"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"io"
	"log"
	"sync"
)

type FeedReader struct {
	database *database.SQLiteDatabase
	feeds    sqlite.Table
	items    sqlite.Table
	search   sqlite.Table
	mu       *sync.Mutex
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

	mu := new(sync.Mutex)

	fr := FeedReader{
		database: db,
		feeds:    f,
		items:    i,
		search:   s,
		mu:       mu,
	}

	return &fr, nil
}

func (fr *FeedReader) AddFeed(feed_url string) (*gofeed.Feed, error) {

	feed, err := fr.ParseFeedURL(feed_url)

	if err != nil {
		return nil, err
	}

	err = fr.IndexFeed(feed)

	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (fr *FeedReader) DumpFeeds(wr io.Writer) error {

	cb := func(f *gofeed.Feed) error {
		wr.Write([]byte(f.FeedLink + "\n"))
		return nil
	}

	err := fr.ListFeedsAll(cb)

	if err != nil {
		log.Fatal()
	}

	return nil
}

func (fr *FeedReader) RefreshFeeds() error {

	fr.mu.Lock()

	defer func() {
		fr.mu.Unlock()
	}()

	// check last update here...

	cb := func(feed *gofeed.Feed) error {

		err := fr.RefreshFeed(feed)

		if err != nil {
			log.Println(feed, err)
		}

		return nil
	}

	return fr.ListFeedsAll(cb)
}

func (fr *FeedReader) Search(q string, opts pagination.PaginatedOptions) (*ItemsResponse, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	// https://www.sqlite.org/fts5.html

	sql := fmt.Sprintf("SELECT feed, guid FROM %s WHERE %s MATCH ? ORDER BY rank", fr.search.Name(), fr.search.Name())

	log.Println("SEARCH", sql, q)

	rsp, err := pagination.QueryPaginated(conn, opts, sql, q)

	if err != nil {
		return nil, err
	}

	guids := make([][]string, 0)

	rows := rsp.Rows()
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

		sql := fmt.Sprintf("SELECT body FROM %s WHERE feed = ? AND guid = ?", fr.items.Name())

		row := conn.QueryRow(sql, feed, guid)
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

func (fr *FeedReader) RemoveFeed(f *gofeed.Feed) error {
	return errors.New("Please write me")
}

func (fr *FeedReader) ListItems(pg_opts pagination.PaginatedOptions) (*ItemsResponse, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	// add "WHERE read=0" toggle
	// add "WHERE feed=..." toggle

	q := fmt.Sprintf("SELECT body FROM %s ORDER BY published ASC, updated ASC", fr.items.Name())

	rsp, err := pagination.QueryPaginated(conn, pg_opts, q)

	if err != nil {
		return nil, err
	}

	items, err := DatabaseRowsToFeedItems(rsp.Rows())

	if err != nil {
		return nil, err
	}

	r := ItemsResponse{
		Items:      items,
		Pagination: rsp.Pagination(),
	}

	return &r, nil
}

func (fr *FeedReader) ListFeedsAll(feed_cb func(f *gofeed.Feed) error) error {

	cb := func(r pagination.PaginatedResponse) error {

		feeds, err := DatabaseRowsToFeeds(r.Rows())

		if err != nil {
			return err
		}

		for _, feed := range feeds {

			err := feed_cb(feed)

			if err != nil {
				log.Println(feed, err)
			}
		}

		return nil
	}

	conn, err := fr.database.Conn()

	if err != nil {
		return err
	}

	sql := fmt.Sprintf("SELECT * FROM %s", fr.feeds.Name())

	opts := pagination.NewDefaultPaginatedOptions()
	opts.PerPage(100)

	return pagination.QueryPaginatedAll(conn, opts, cb, sql)
}

func (fr *FeedReader) ListFeeds(pg_opts pagination.PaginatedOptions) (*FeedsResponse, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT body FROM %s ORDER BY updated ASC", fr.feeds.Name())

	rsp, err := pagination.QueryPaginated(conn, pg_opts, q)

	if err != nil {
		return nil, err
	}

	rows := rsp.Rows()
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

func (fr *FeedReader) RefreshFeed(feed *gofeed.Feed) error {

	f2, err := fr.ParseFeedURL(feed.FeedLink)

	if err != nil {
		return err
	}

	err = fr.IndexFeed(f2)

	if err != nil {
		return err
	}

	return nil
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

func (fr *FeedReader) IndexFeed(feed *gofeed.Feed) error {

	items := feed.Items
	feed.Items = nil

	err := fr.feeds.IndexRecord(fr.database, feed)

	if err != nil {
		return err
	}

	for _, item := range items {

		rec := tables.ItemsRecord{
			Feed: feed,
			Item: item,
		}

		err = fr.items.IndexRecord(fr.database, &rec)

		if err != nil {
			return err
		}

		err = fr.search.IndexRecord(fr.database, &rec)

		if err != nil {
			return err
		}
	}

	return nil
}
