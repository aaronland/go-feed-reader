package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/http"
	"log"
	gohttp "net/http"
	"os"
	"time"
)

func main() {

	var host = flag.String("host", "localhost", "")
	var port = flag.Int("port", 8080, "")

	var dsn = flag.String("dsn", ":memory:", "")
	var refresh = flag.Int("refresh", 15, "")

	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	go fr.RefreshFeeds()

	feeds_handler, err := http.FeedsHandler(fr)

	if err != nil {
		log.Fatal(err)
	}

	recent_handler, err := http.RecentItemsHandler(fr)

	if err != nil {
		log.Fatal(err)
	}

	item_handler, err := http.ItemHandler(fr)

	if err != nil {
		log.Fatal(err)
	}

	add_handler, err := http.AddHandler(fr)

	if err != nil {
		log.Fatal(err)
	}

	search_handler, err := http.SearchHandler(fr)

	if err != nil {
		log.Fatal(err)
	}

	go func() {

		d := time.Duration(*refresh)

		ticker := time.NewTicker(time.Minute * d)

		for range ticker.C {

			log.Println("refresh feeds")

			err := fr.RefreshFeeds()

			if err != nil {
				log.Println(err)
			}
		}

	}()

	mux := gohttp.NewServeMux()

	mux.Handle("/", feeds_handler)
	mux.Handle("/feeds", feeds_handler)
	mux.Handle("/search", search_handler)
	mux.Handle("/item", item_handler)
	mux.Handle("/add", add_handler)
	mux.Handle("/recent", recent_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening on %s\n", endpoint)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
