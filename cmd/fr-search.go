package main

import (
	"flag"
	"github.com/aaronland/go-feed-reader"
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

	items, err := fr.Search(*q)

	if err != nil {
		log.Fatal(err)
	}

	for _, i := range items {
		log.Println(i.Title, i.Link)
	}

	os.Exit(0)
}
