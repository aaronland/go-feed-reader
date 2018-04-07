package tables

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	"time"
)

type ItemsTable struct {
	sqlite.Table
	name string
}

// I don't like this...

type ItemsRecord struct {
	Feed *gofeed.Feed
	Item *gofeed.Item
}

func NewItemsTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewItemsTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewItemsTable() (sqlite.Table, error) {

	t := ItemsTable{
		name: "items",
	}

	return &t, nil
}

func (t *ItemsTable) Name() string {
	return t.name
}

func (t *ItemsTable) Schema() string {

	sql := `CREATE TABLE %s (
	    	feed TEXT NOT NULL,
		guid TEXT NOT NULL,
		link TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		body JSON NOT NULL,
		published INTEGER,
		updated INTEGER,
		read INTEGER,
		saved INTEGER
	);

	CREATE UNIQUE INDEX %s_by_guid ON %s (feed, guid);
	CREATE INDEX %s_by_published ON %s (published, updated);
	CREATE INDEX %s_by_read ON %s (read, published, updated);	
	CREATE INDEX %s_by_feed ON %s (feed_link, read, published, updated);
	`

	return fmt.Sprintf(sql, t.Name(), t.Name(), t.Name(), t.Name(), t.Name(), t.Name(), t.Name(), t.Name(), t.Name())
}

func (t *ItemsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *ItemsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	rec := i.(*ItemsRecord) // I don't like this...
	return t.IndexItem(db, rec.Feed, rec.Item)
}

func (t *ItemsTable) IndexItem(db sqlite.Database, f *gofeed.Feed, i *gofeed.Item) error {

	b, err := json.Marshal(i)

	if err != nil {
		return err
	}

	body := string(b)

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		feed, guid, link, title, description, body, published, updated
	) VALUES (
	  	 ?, ?, ?, ?, ?, ?, ?, ?
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

	_, err = stmt.Exec(f.Link, i.GUID, i.Link, i.Title, i.Description, body, published_ts, updated_ts)

	if err != nil {
		return err
	}

	return tx.Commit()
}
