package reader

// this will eventually be moved in to go-whosonfirst-sqlite or something
// (20180407/thisisaaronland)

// also this does _not_ work yet...
// (20180407/thisisaaronland)

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Paginated struct {
	Rows       *sql.Rows
	TotalCount int
}

func QueryPaginated(db *sql.DB, query string, args ...interface{}) (*Paginated, error) {

	done_ch := make(chan bool)
	err_ch := make(chan error)
	count_ch := make(chan int)
	rows_ch := make(chan *sql.Rows)

	go func() {

		defer func() {
			done_ch <- true
		}()

		parts := strings.Split(query, " FROM ")
		parts = strings.Split(parts[1], " LIMIT ")

		conditions := parts[0]

		count_query := fmt.Sprintf("SELECT COUNT(*) FROM %s", conditions)

		row := db.QueryRow(count_query)

		var count int
		err := row.Scan(&count)

		if err != nil {
			err_ch <- err
			return
		}

		count_ch <- count
	}()

	go func() {

		defer func() {
			done_ch <- true
		}()

		rows, err := db.Query(query, args...)

		if err != nil {
			err_ch <- err
			return
		}

		rows_ch <- rows
	}()

	var total_count int
	var rows *sql.Rows

	remaining := 2

	for remaining > 0 {

		select {
		case <-done_ch:
			remaining -= 1
		case e := <-err_ch:
			return nil, e
		case i := <-count_ch:
			total_count = i
		case r := <-rows_ch:
			rows = r
		default:
			//
		}
	}

	pg := Paginated{
		TotalCount: total_count,
		Rows:       rows,
	}

	return &pg, errors.New("Please finish me")
}
