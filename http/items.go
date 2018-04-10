package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/aaronland/go-sql-pagination"	
	"github.com/arschles/go-bindata-html-template"
	"github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	gohttp "net/http"
	"strconv"
)

type ItemsVars struct {
	PageTitle  string
	Items      []*gofeed.Item
	Pagination pagination.Pagination
}

func ItemsHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_items.html",
		"templates/html/inc_pagination.html",
		"templates/html/inc_foot.html",
	}

	funcs := template.FuncMap{
		"strip_tags": strip.StripTags,
	}

	t, err := template.New("items", html.Asset).Funcs(funcs).ParseFiles(files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		opts := pagination.NewDefaultPaginatedOptions()

		query := req.URL.Query()
		str_page := query.Get("page")

		if str_page != "" {

			page, err := strconv.Atoi(str_page)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			opts.Page(page)
		}

		q_rsp, err := fr.ListItems(opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := ItemsVars{
			PageTitle:  "Recent items",
			Items:      q_rsp.Items,
			Pagination: q_rsp.Pagination,
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
