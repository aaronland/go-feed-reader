package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-sanitize"
	_ "log"
	gohttp "net/http"
)

type ResultsVars struct {
	Items []*gofeed.Item
	Query string
}

func SearchHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	t, err := CompileTemplates([]string{
		"templates/html/inc_head.html",
		"templates/html/inc_foot.html",
		"templates/html/inc_items.html",
		"templates/html/inc_search_form.html",		
		"templates/html/search_query.html",
		"templates/html/search_results.html",
	})

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		query := req.URL.Query()
		raw_q := query.Get("q")

		opts := sanitize.DefaultOptions()
		q, err := sanitize.SanitizeString(raw_q, opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		if q == "" {

			err := t.ExecuteTemplate(rsp, "search_query.html", "")

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return
		}

		items, err := fr.Search(q)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := ResultsVars{
			Items: items,
			Query: q,
		}

		err = t.ExecuteTemplate(rsp, "search_results.html", vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
