package reader

import (
       "errors"
	"github.com/aaronland/go-feed-reader/tables"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
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
		items:	  i,
	}

	return &fr, nil
}

func (fr *FeedReader) AddFeed(f *gofeed.Feed) error {
     return fr.feeds.IndexRecord(fr.database, f)
}

func (fr *FeedReader) RemoveFeed(f *gofeed.Feed) error {
     return errors.New("Please write me")
}

func (fr *FeedReader) ListFeeds(f *gofeed.Feed) ([]*gofeed.Feed, error) {
     return nil, errors.New("Please write me")
}

func (fr *FeedReader) RefreshFeeds(feeds []*gofeed.Feed) error {
     return errors.New("Please write me")
}

func (fr *FeedReader) RefreshFeed(f *gofeed.Feed) error {
     return errors.New("Please write me")
}


