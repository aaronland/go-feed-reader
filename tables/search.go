package tables

// https://www.sqlite.org/fts5.html

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	_ "log"
)

type SearchTable struct {
	sqlite.Table
	name string
}

func NewSearchTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewSearchTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewSearchTable() (sqlite.Table, error) {

	t := SearchTable{
		name: "search",
	}

	return &t, nil
}

func (t *SearchTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *SearchTable) Name() string {
	return t.name
}

func (t *SearchTable) Schema() string {

	schema := `CREATE VIRTUAL TABLE %s USING fts5(
		feed, guid,
		title, link, content 
	);`

	// so dumb...
	return fmt.Sprintf(schema, t.Name())
}

func (t *SearchTable) IndexRecord(db sqlite.Database, i interface{}) error {
	rec := i.(*ItemRecord) // I don't like this...
	return t.IndexItem(db, rec.Feed, rec.Item)
}

func (t *SearchTable) IndexItem(db sqlite.Database, f *gofeed.Feed, i *gofeed.Item) error {

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		feed, guid,
		title, link, content		      
		) VALUES (
		?, ?,
		?, ?, ?
		)`, t.Name()) // ON CONFLICT DO BLAH BLAH BLAH

	args := []interface{}{
		f.Link, i.GUID,
		i.Title, i.Link, i.Content,
	}

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	s, err := tx.Prepare(fmt.Sprintf("DELETE FROM %s WHERE feed = ? AND guid = ?", t.Name()))

	if err != nil {
		return err
	}

	defer s.Close()

	_, err = s.Exec(f.Link, i.GUID)

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(args...)

	if err != nil {
		return err
	}

	return tx.Commit()
}
