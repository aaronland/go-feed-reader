package main

import (
	"flag"
	"github.com/aaronland/go-feed-reader"
	"log"
	_ "net/http"
)

func main() {

	var dsn = flag.String("dsn", ":memory:", "")
	var feed_url = flag.String("feed", "", "")

	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	feed, err := fr.ParseFeedURL(*feed_url)

	if err != nil {
		log.Fatal(err)
	}

	err = fr.IndexFeed(feed)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")
}
