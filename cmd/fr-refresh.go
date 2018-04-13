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

	err = fr.RefreshFeeds()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
