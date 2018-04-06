package tables

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
)

type ItemsTable struct {
	sqlite.Table
	name string
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
		name: "feeds",
	}

	return &t, nil
}

func (t *ItemsTable) Name() string {
	return t.name
}

func (t *ItemsTable) Schema() string {

	// feed TEXT NOT NULL,

	sql := `CREATE TABLE %s (
		guid TEXT NOT NULL PRIMARY KEY,
		link TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		body JSON NOT NULL
	);`

	// lastmodified INTEGER

	return fmt.Sprintf(sql, t.Name())
}

func (t *ItemsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *ItemsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	return t.IndexItem(db, i.(*gofeed.Item))
}

func (t *ItemsTable) IndexItem(db sqlite.Database, i *gofeed.Item) error {

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
		guid, link, title, description, body
	) VALUES (
	  	 ?, ?, ?, ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(i.GUID, i.Link, i.Title, i.Description, body)

	if err != nil {
		return err
	}

	return tx.Commit()
}
