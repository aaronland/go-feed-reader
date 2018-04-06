package tables

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
)

type FeedsTable struct {
	sqlite.Table
	name string
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
		title TEXT NOT NULL,
		author TEXT,
		link TEXT NOT NULL PRIMARY KEY,
		body JSON NOT NULL,
		lastmodified INTEGER
	);`

	return fmt.Sprintf(sql, t.Name())
}

func (t *FeedsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *FeedsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	return t.IndexFeed(db, i.(*gofeed.Feed))
}

func (t *FeedsTable) IndexFeed(db sqlite.Database, f *gofeed.Feed) error {

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

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		title, author, link, body, lastmodified
	) VALUES (
	  	 ?, ?, ?, ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(f.Title, f.Author.Name, f.Link, str_body, f.Updated)

	if err != nil {
		return err
	}

	return tx.Commit()
}
