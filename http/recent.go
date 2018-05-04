package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/aaronland/go-sql-pagination"
	"github.com/mmcdole/gofeed"
	_ "log"
	gohttp "net/http"
	"net/url"
	"strconv"
)

type RecentItemsVars struct {
	PageTitle  string
	Items      []*gofeed.Item
	Pagination pagination.Pagination
	URL        *url.URL
	User       user.User
}

func RecentItemsHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_pagination.html",
		"templates/html/inc_items.html",
		"templates/html/inc_pagination.html",
		"templates/html/inc_foot.html",
	}

	t, err := CompileTemplate("items", files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		u := EnsureLoggedIn(fr, rsp, req)

		if u == nil {
			return
		}

		str_page, err := GetString(req, "page")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		str_feed, err := GetString(req, "feed")

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		pg_opts := pagination.NewDefaultPaginatedOptions()

		if str_page != "" {

			page, err := strconv.Atoi(str_page)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			pg_opts.Page(page)
		}

		ls_opts := reader.NewDefaultListItemsOptions()

		if str_feed != "" {

			u, err := url.Parse(str_feed)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			ls_opts.FeedURL = u.String()
		}

		q_rsp, err := fr.ListItemsForUser(u, ls_opts, pg_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := RecentItemsVars{
			PageTitle:  "Recent items",
			Items:      q_rsp.Items,
			Pagination: q_rsp.Pagination,
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
