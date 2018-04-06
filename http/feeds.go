package http

import (
	"fmt"
	"github.com/aaronland/go-feed-reader"
	"github.com/mmcdole/gofeed"
	"html/template"
	gohttp "net/http"
)

type HTMLVars struct {
	Feeds []*gofeed.Feed
}

func FeedsHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	t, err := template.New("feeds").Parse(``)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		feeds, err := fr.ListFeeds()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := HTMLVars{
			Feeds: feeds,
		}

		err = t.ExecuteTemplate(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
