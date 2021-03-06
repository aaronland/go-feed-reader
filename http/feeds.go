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

type FeedVars struct {
	PageTitle  string
	Feeds      []*gofeed.Feed
	Pagination pagination.Pagination
	URL        *url.URL
	User       user.User
}

func FeedsHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_pagination.html",
		"templates/html/inc_feeds.html",
		"templates/html/inc_pagination.html",
		"templates/html/inc_foot.html",
	}

	t, err := CompileTemplate("feeds", files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		u := EnsureLoggedIn(fr, rsp, req)

		if u == nil {
			return
		}

		pg_opts := pagination.NewDefaultPaginatedOptions()

		query := req.URL.Query()
		str_page := query.Get("page")

		if str_page != "" {

			page, err := strconv.Atoi(str_page)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			pg_opts.Page(page)
		}

		results, err := fr.ListFeedsForUser(u, pg_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := FeedVars{
			PageTitle:  "",
			Feeds:      results.Feeds,
			Pagination: results.Pagination,
			URL:        req.URL,
			User:       u,
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
