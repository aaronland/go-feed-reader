package tables

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	_ "log"
	"time"
)

type UserFeedsTable struct {
	sqlite.Table
	name string
}

type UserFeed struct {
	Feed *gofeed.Feed
	User user.User
}

func NewUserFeedsTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewUserFeedsTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewUserFeedsTable() (sqlite.Table, error) {

	t := UserFeedsTable{
		name: "user_feeds",
	}

	return &t, nil
}

func (t *UserFeedsTable) Name() string {
	return t.name
}

func (t *UserFeedsTable) Schema() string {

	sql := `CREATE TABLE %s (
	    	id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		feed_link TEXT NULL,		     
	);

	CREATE UNIQUE INDEX %s_unq_feed ON %s (user_id, feed_link);
	`

	return fmt.Sprintf(sql, t.Name(), t.Name())
}

func (t *UserFeedsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *UserFeedsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	rec := i.(*UserFeed)
	return t.IndexFeed(db, rec)
}

func (t *UserFeedsTable) IndexUserFeed(db sqlite.Database, f *UserFeed) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		user_id, feed_link
	) VALUES (
	  	 ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(f.User.Id(), f.FeedLink)

	if err != nil {
		return err
	}

	return tx.Commit()
}
