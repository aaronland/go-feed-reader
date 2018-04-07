package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/http"
	"log"
	gohttp "net/http"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "")
	var port = flag.Int("port", 8080, "")

	var dsn = flag.String("dsn", ":memory:", "")

	flag.Parse()

	fr, err := reader.NewFeedReader(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	feeds_handler, err := http.FeedsHandler(fr)

	if err != nil {
		log.Fatal(err)
	}

	mux := gohttp.NewServeMux()

	mux.Handle("/feeds", feeds_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
