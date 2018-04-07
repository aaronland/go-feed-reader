package main

import (
	"flag"
	"github.com/aaronland/go-feed-reader"
	"log"
	"os"
	"time"
)

func main() {

	var dsn = flag.String("dsn", ":memory:", "")
	var daemon = flag.Bool("daemon", false, "")

	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	for {

		feeds, err := fr.ListFeeds()

		if err != nil {
			log.Fatal(err)
		}

		err = fr.RefreshFeeds(feeds)

		if err != nil {
			log.Fatal(err)
		}

		if !*daemon {
			break
		}

		time.Sleep(time.Minute * 30)
	}

	os.Exit(0)
}
