package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronland/go-feed-reader/tables"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	_ "log"
)

type FeedReader struct {
	database *database.SQLiteDatabase
	feeds    sqlite.Table
	items    sqlite.Table
	search   sqlite.Table
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

	fr := FeedReader{
		database: db,
		feeds:    f,
		items:    i,
		search:   s,
	}

	return &fr, nil
}

func (fr *FeedReader) Search(q string) ([]*gofeed.Item, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	// https://www.sqlite.org/fts5.html

	sql := fmt.Sprintf("SELECT feed, guid FROM %s(?) ORDER BY rank", fr.search.Name())

	rows, err := conn.Query(sql, q)

	if err != nil {
		return nil, err
	}

	guids := make([][]string, 0)

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

	return items, nil
}

func (fr *FeedReader) RemoveFeed(f *gofeed.Feed) error {
	return errors.New("Please write me")
}

func (fr *FeedReader) ListItems() ([]*gofeed.Item, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	// add "WHERE read=0" toggle

	q := fmt.Sprintf("SELECT body FROM %s ORDER BY published ASC, updated ASC", fr.items.Name())

	opts := NewDefaultPaginationOptions()

	rsp, err := QueryPaginated(conn, opts, q)

	if err != nil {
		return nil, err
	}

	/*
		rows, err := conn.Query(q)

		if err != nil {
			return nil, err
		}
	*/

	items, err := DatabaseRowsToFeedItems(rsp.Rows())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (fr *FeedReader) ListFeeds() ([]*gofeed.Feed, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT body FROM %s ORDER BY updated ASC", fr.feeds.Name())

	rows, err := conn.Query(q)

	if err != nil {
		return nil, err
	}

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

	return feeds, nil
}

func (fr *FeedReader) RefreshFeeds(feeds []*gofeed.Feed) error {

	for _, f := range feeds {

		f2, err := fr.RefreshFeed(f)

		if err != nil {
			return err
		}

		err = fr.IndexFeed(f2)

		if err != nil {
			return err
		}
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

func (fr *FeedReader) RefreshFeed(feed *gofeed.Feed) (*gofeed.Feed, error) {

	return fr.ParseFeedURL(feed.FeedLink)
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
