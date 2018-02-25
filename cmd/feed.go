package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-feed-reader"
	"github.com/mmcdole/gofeed"
	"log"
	_ "net/http"
)

func main() {

	var dsn = flag.String("dsn", ":memory:", "")

	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(fr)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://whosonfirst.org/blog/rss_20.xml")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(feed.Title)

	err = fr.AddFeed(feed)

	if err != nil {
		log.Fatal(err)
	}

}
