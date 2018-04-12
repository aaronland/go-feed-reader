package main

import (
	"flag"
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-sql-pagination"
	"log"
	"os"
)

func main() {

	var dsn = flag.String("dsn", ":memory:", "")
	var q = flag.String("query", "", "")

	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	pg_opts := pagination.NewDefaultPaginatedOptions()

	rsp, err := fr.Search(*q, pg_opts)

	if err != nil {
		log.Fatal(err)
	}

	for _, i := range rsp.Items {
		log.Println(i.Title, i.Link)
	}

	os.Exit(0)
}
