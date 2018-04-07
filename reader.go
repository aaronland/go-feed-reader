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

	fr := FeedReader{
		database: db,
		feeds:    f,
		items:    i,
	}

	return &fr, nil
}

func (fr *FeedReader) RemoveFeed(f *gofeed.Feed) error {
	return errors.New("Please write me")
}

func (fr *FeedReader) ListItems() ([]*gofeed.Item, error) {
	return nil, errors.New("Please write me")
}

func (fr *FeedReader) ListFeeds() ([]*gofeed.Feed, error) {

	conn, err := fr.database.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT body FROM %s", fr.feeds.Name())

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

func (fr *FeedReader) UpdateFeeds(feeds []*gofeed.Feed) error {

	// please make me a generator or equivalent
	// (20180405/thisisaaronland)

	feeds, err := fr.ListFeeds()

	if err != nil {
		return err
	}

	for _, f := range feeds {

		f, err = fr.RefreshFeed(f)

		if err != nil {
			return err
		}

		err = fr.IndexFeed(f)

		if err != nil {
			return err
		}
	}

	return nil
}

func (fr *FeedReader) ParseFeedURL(feed_url string) (*gofeed.Feed, error) {

	fp := gofeed.NewParser()
	return fp.ParseURL(feed_url)
}

func (fr *FeedReader) RefreshFeed(feed *gofeed.Feed) (*gofeed.Feed, error) {

	return fr.ParseFeedURL(feed.FeedLink)
}

func (fr *FeedReader) IndexFeed(feed *gofeed.Feed) error {

	err := fr.feeds.IndexRecord(fr.database, feed)

	if err != nil {
		return err
	}

	for _, item := range feed.Items {

		rec := tables.ItemsRecord{
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
