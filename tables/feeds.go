package tables

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	_ "log"
	"time"
)

type FeedsTable struct {
	sqlite.Table
	name string
}

type FeedRecord struct {
	Feed *gofeed.Feed
}

func NewFeedsTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewFeedsTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewFeedsTable() (sqlite.Table, error) {

	t := FeedsTable{
		name: "feeds",
	}

	return &t, nil
}

func (t *FeedsTable) Name() string {
	return t.name
}

func (t *FeedsTable) Schema() string {

	sql := `CREATE TABLE %s (
	    	id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		link TEXT NOT NULL,
		feed_link TEXT NULL,		     
		body JSON NOT NULL,
		published INTEGER,		     
		updated INTEGER
	);

	CREATE INDEX %s_by_published ON %s (published, updated);
	`

	return fmt.Sprintf(sql, t.Name(), t.Name(), t.Name())
}

func (t *FeedsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *FeedsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	rec := i.(*FeedRecord) // I don't like this...
	return t.IndexFeed(db, rec.Feed)
}

func (t *FeedsTable) IndexFeed(db sqlite.Database, f *gofeed.Feed) error {

	if f.Title == "" {
		return errors.New("Unable to determine feed title")
	}

	if f.Link == "" {
		return errors.New("Unable to determine feed link")
	}

	body, err := json.Marshal(f)

	if err != nil {
		return err
	}

	str_body := string(body)

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		title, link, feed_link, body, published, updated
	) VALUES (
	  	 ?, ?, ?, ?, ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	var published_ts int64
	var updated_ts int64

	published := f.PublishedParsed
	updated := f.UpdatedParsed

	if published == nil {
		now := time.Now()
		published_ts = now.Unix()
	} else {
		published_ts = published.Unix()
	}

	if updated == nil {
		now := time.Now()
		updated_ts = now.Unix()
	} else {
		updated_ts = updated.Unix()
	}

	_, err = stmt.Exec(f.Title, f.Link, f.FeedLink, str_body, published_ts, updated_ts)

	if err != nil {
		return err
	}

	return tx.Commit()
}
