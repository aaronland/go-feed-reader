package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/arschles/go-bindata-html-template"
	"github.com/mmcdole/gofeed"
	gohttp "net/http"
)

type FeedVars struct {
     PageTitle string
	Feeds []*gofeed.Feed
}

func FeedsHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_feeds.html",
		"templates/html/inc_foot.html",
	}

	t, err := template.New("feeds", html.Asset).ParseFiles(files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		feeds, err := fr.ListFeeds()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := FeedVars{
		     	PageTitle: "",
			Feeds: feeds,
		}

		err = t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
