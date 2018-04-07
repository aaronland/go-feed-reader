package main

import (
	"flag"
	"github.com/aaronland/go-feed-reader"
	"log"
	"os"
)

func main() {

	var dsn = flag.String("dsn", ":memory:", "")

	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	feeds, err := fr.ListFeeds()

	if err != nil {
		log.Fatal()
	}

	for _, f := range feeds {
		os.Stdout.Write([]byte(f.FeedLink + "\n"))
	}

	os.Exit(0)
}
