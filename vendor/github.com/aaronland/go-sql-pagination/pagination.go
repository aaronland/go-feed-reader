package pagination

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
)

type PaginatedOptions interface {
	PerPage(...int) int
	Page(...int) int
	Spill(...int) int
	Column(...string) string
}

type PaginatedResponse interface {
	Rows() *sql.Rows
	Pagination() Pagination
}

type Pagination interface {
	Total() int
	PerPage() int
	Page() int
	Pages() int
}

type DefaultPaginatedResponse struct {
	rows       *sql.Rows
	pagination Pagination
}

func (r *DefaultPaginatedResponse) Rows() *sql.Rows {
	return r.rows
}

func (r *DefaultPaginatedResponse) Pagination() Pagination {
	return r.pagination
}

type DefaultPagination struct {
	Pagination
	total    int
	per_page int
	page     int
	pages    int
}

func (p *DefaultPagination) Total() int {
	return p.total
}

func (p *DefaultPagination) PerPage() int {
	return p.per_page
}

func (p *DefaultPagination) Page() int {
	return p.page
}

func (p *DefaultPagination) Pages() int {
	return p.pages
}

type DefaultPaginatedOptions struct {
	PaginatedOptions
	per_page int
	page     int
	spill    int
	column   string
}

func (o *DefaultPaginatedOptions) PerPage(args ...int) int {

	if len(args) == 1 {
		o.per_page = args[0]
	}
	return o.per_page
}

func (o *DefaultPaginatedOptions) Page(args ...int) int {

	if len(args) == 1 {
		o.page = args[0]
	}

	return o.page
}

func (o *DefaultPaginatedOptions) Spill(args ...int) int {

	if len(args) == 1 {
		o.spill = args[0]
	}

	return o.spill
}

func (o *DefaultPaginatedOptions) Column(args ...string) string {

	if len(args) == 1 {
		o.column = args[0]
	}

	return o.column
}

func NewDefaultPaginatedOptions() PaginatedOptions {

	opts := DefaultPaginatedOptions{
		per_page: 10,
		page:     1,
		spill:    2,
		column:   "*",
	}

	return &opts
}

func QueryPaginated(db *sql.DB, opts PaginatedOptions, query string, args ...interface{}) (PaginatedResponse, error) {

	done_ch := make(chan bool)
	err_ch := make(chan error)
	count_ch := make(chan int)
	rows_ch := make(chan *sql.Rows)

	var page int
	var per_page int
	var spill int

	go func() {

		defer func() {
			done_ch <- true
		}()

		parts := strings.Split(query, " FROM ")
		parts = strings.Split(parts[1], " LIMIT ")

		conditions := parts[0]

		count_query := fmt.Sprintf("SELECT COUNT(%s) FROM %s", opts.Column(), conditions)
		log.Println("COUNT", count_query)

		row := db.QueryRow(count_query)

		var count int
		err := row.Scan(&count)

		if err != nil {
			err_ch <- err
			return
		}

		log.Println("COUNT", count)
		count_ch <- count
	}()

	go func() {

		defer func() {
			done_ch <- true
		}()

		// please make fewer ((((())))) s
		// (20180409/thisisaaronland)

		page = int(math.Max(1.0, float64(opts.Page())))
		per_page = int(math.Max(1.0, float64(opts.PerPage())))
		spill = int(math.Max(1.0, float64(opts.Spill())))

		if spill >= per_page {
			spill = per_page - 1
		}

		offset := 0
		limit := per_page

		offset = (page - 1) * per_page

		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
		log.Println("QUERY", query)

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

	pages := int(math.Ceil(float64(total_count) / float64(per_page)))
	log.Println("PAGES", pages)

	pg := DefaultPagination{
		total:    total_count,
		per_page: per_page,
		page:     page,
		pages:    pages,
	}

	rsp := DefaultPaginatedResponse{
		pagination: &pg,
		rows:       rows,
	}

	return &rsp, nil
}
