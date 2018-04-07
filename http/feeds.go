package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/mmcdole/gofeed"
	"html/template"
	gohttp "net/http"
)

type FeedVars struct {
	Feeds []*gofeed.Feed
}

func FeedsHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	tpl, err := html.Asset("templates/html/feeds.html")

	if err != nil {
		return nil, err
	}

	t, err := template.New("feeds").Parse(string(tpl))

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
			Feeds: feeds,
		}

		err = t.ExecuteTemplate(rsp, "feeds", vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
