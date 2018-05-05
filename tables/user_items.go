package tables

import (
	"fmt"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	_ "log"
)

type UserItemsTable struct {
	sqlite.Table
	name string
}

type UserItem struct {
	Feed *gofeed.Feed
	Item *gofeed.Item
	User user.User
}

func NewUserItemsTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewUserItemsTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewUserItemsTable() (sqlite.Table, error) {

	t := UserItemsTable{
		name: "user_items",
	}

	return &t, nil
}

func (t *UserItemsTable) Name() string {
	return t.name
}

func (t *UserItemsTable) Schema() string {

	sql := `CREATE TABLE %s (
	    	id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		feed_link TEXT NULL,
		item_guid TEXT NULL,
	);

	CREATE UNIQUE INDEX %s_unq_item ON %s (user_id, feed_link, item_guid);
	`

	return fmt.Sprintf(sql, t.Name(), t.Name())
}

func (t *UserItemsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *UserItemsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	rec := i.(*UserItem)
	return t.IndexUserItem(db, rec)
}

func (t *UserItemsTable) IndexUserItem(db sqlite.Database, f *UserItem) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		user_id, feed_link, item_guid
	) VALUES (
	  	 ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(f.User.Id(), f.Feed.Link, f.Item.GUID)

	if err != nil {
		return err
	}

	return tx.Commit()
}
