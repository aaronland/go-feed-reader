package reader

// this will eventually be moved in to go-whosonfirst-sqlite or something
// (20180407/thisisaaronland)

// also this does _not_ work yet...
// (20180407/thisisaaronland)

import (
	"database/sql"
	"fmt"
	_ "log"		
	"strings"
)

type PaginationOptions struct {
     PerPage int
     Page int
     Column string
}

type PaginatedRows struct {
	Rows       *sql.Rows
	Total int
     	PerPage int
     	Page int
	Pages int	     
}

func DefaultPaginationOptions() *PaginationOptions {

     opts := PaginationOptions {
     	  PerPage: 10,
     	  Page: 1,
	  Column: "*",
     }

     return &opts
}

func QueryPaginated(db *sql.DB, opts *PaginationOptions, query string, args ...interface{}) (*PaginatedRows, error) {

	done_ch := make(chan bool)
	err_ch := make(chan error)
	count_ch := make(chan int)
	rows_ch := make(chan *sql.Rows)

	go func() {

		defer func() {
			done_ch <- true
		}()

		parts := strings.Split(query, " FROM ")
		conditions := parts[1]

		count_query := fmt.Sprintf("SELECT COUNT(%s) FROM %s", opts.Column, conditions)
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

		offset := 0
		limit := opts.PerPage

		if opts.Page > 1 {
			offset = (opts.Page - 1) * opts.PerPage
		}
		
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)		
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

	pages := 0	// FIX ME
	
	pg := PaginatedRows{
		Total: total_count,
		PerPage: opts.PerPage,
		Page: opts.Page,
		Pages: pages,
		Rows:       rows,
	}

	return &pg, nil
}
