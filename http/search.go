package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/aaronland/go-sql-pagination"
	"github.com/arschles/go-bindata-html-template"
	"github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	_ "log"
	gohttp "net/http"
	"net/url"
	"strconv"
)

type SearchFormVars struct {
	PageTitle string
}

type ResultsVars struct {
	PageTitle  string
	Items      []*gofeed.Item
	Query      string
	Pagination pagination.Pagination
	URL        *url.URL
}

func SearchHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	query_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_foot.html",
	}

	results_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_search_results.html",
		"templates/html/inc_pagination.html",
		"templates/html/inc_items.html",
		"templates/html/inc_pagination.html",
		"templates/html/inc_foot.html",
	}

	results_funcs := template.FuncMap{
		"strip_tags": strip.StripTags,
	}

	query_t, err := template.New("query", html.Asset).ParseFiles(query_files...)

	if err != nil {
		return nil, err
	}

	results_t, err := template.New("results", html.Asset).Funcs(results_funcs).ParseFiles(results_files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		pg_opts := pagination.NewDefaultPaginatedOptions()
		pg_opts.Column("feed")

		q, err := GetString(req, "q")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		p, err := GetString(req, "page")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		if q == "" {

			vars := SearchFormVars{
				PageTitle: "",
			}

			err := query_t.ExecuteTemplate(rsp, "query", vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return
		}

		if p != "" {

			page, err := strconv.Atoi(p)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			pg_opts.Page(page)
		}

		results, err := fr.Search(q, pg_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := ResultsVars{
			PageTitle:  "",
			Items:      results.Items,
			Pagination: results.Pagination,
			URL:        req.URL,
			Query:      q,
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
