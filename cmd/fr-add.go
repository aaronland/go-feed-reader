package main

import (
	"flag"
	"github.com/aaronland/go-feed-reader"
	"log"
)

func main() {

	var dsn = flag.String("dsn", ":memory:", "")
	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	for _, url := range flag.Args() {

		feed, err := fr.ParseFeedURL(url)

		if err != nil {
			log.Fatal(err)
		}

		err = fr.IndexFeed(feed)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("OK")
}
