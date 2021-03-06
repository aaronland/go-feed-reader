package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/aaronland/go-sql-pagination"
	"github.com/mmcdole/gofeed"
	gohttp "net/http"
	"net/url"
	"strconv"
)

type SearchFormVars struct {
	PageTitle string
	User      user.User
}

type ResultsVars struct {
	PageTitle  string
	Items      []*gofeed.Item
	Query      string
	Pagination pagination.Pagination
	URL        *url.URL
	User       user.User
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

	query_t, err := CompileTemplate("query", query_files...)

	if err != nil {
		return nil, err
	}

	results_t, err := CompileTemplate("results", results_files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		u := EnsureLoggedIn(fr, rsp, req)

		if u == nil {
			return
		}

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
				User:      u,
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

		results, err := fr.SearchForUser(u, q, pg_opts)

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
			User:       u,
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
