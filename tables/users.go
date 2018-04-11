package tables

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
)

type UsersTable struct {
	sqlite.Table
	name string
}

type User struct {
     	Id int64
	Name string
	Password string
	Salt string		 
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
		name: "items",
	}

	return &t, nil
}

func (t *UsersTable) Name() string {
	return t.name
}

func (t *UsersTable) Schema() string {

	sql := `CREATE TABLE %s (
	    	id INTEGER PRIMARY KEY AUTO_INCREMENT,
		name TEXT NOT NULL,
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
	user := i.(*User)
	return t.IndexUser(db, user)
}

func (t *UsersTable) IndexUser(db sqlite.Database, u *User) error {

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		name, password, salt
	) VALUES (
	  	 ?, ?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Password, u.Salt)

	if err != nil {
		return err
	}

	return tx.Commit()
}