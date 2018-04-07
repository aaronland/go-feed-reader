package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/arschles/go-bindata-html-template"
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

	query_files := []string{
		"templates/html/inc_head.html",		
		"templates/html/inc_search_form.html",
		"templates/html/inc_foot.html",
	}

	results_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_search_form.html",
		"templates/html/inc_search_results.html",		
		"templates/html/inc_items.html",		
		"templates/html/inc_foot.html",
	}

	query_t, err := template.New("query", html.Asset).ParseFiles(query_files...)

	if err != nil {
		return nil, err
	}

	results_t, err := template.New("results", html.Asset).ParseFiles(results_files...)

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

			err := query_t.ExecuteTemplate(rsp, "query", "")

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

		err = results_t.ExecuteTemplate(rsp, "results", vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
