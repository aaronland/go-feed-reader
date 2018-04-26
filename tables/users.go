package tables

import (
	"fmt"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
)

type UsersTable struct {
	sqlite.Table
	name string
}

func NewUsersTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewUsersTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewUsersTable() (sqlite.Table, error) {

	t := UsersTable{
		name: "users",
	}

	return &t, nil
}

func (t *UsersTable) Name() string {
	return t.name
}

func (t *UsersTable) Schema() string {

	sql := `CREATE TABLE %s (
	    	id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,		
		password TEXT NOT NULL,
		salt TEXT NOT NULL
	);

	`

	return fmt.Sprintf(sql, t.Name())
}

func (t *UsersTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *UsersTable) IndexRecord(db sqlite.Database, i interface{}) error {
	user := i.(user.User)
	return t.IndexUser(db, user)
}

func (t *UsersTable) IndexUser(db sqlite.Database, u user.User) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		name, email, password, salt
	) VALUES (
	  	 ?, ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	pswd := u.Password()

	_, err = stmt.Exec(u.Username(), u.Email(), pswd.Digest(), "SALT")

	if err != nil {
		return err
	}

	return tx.Commit()
}
