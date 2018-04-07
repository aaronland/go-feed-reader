package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-sanitize"	
	"html/template"
	_ "log"		
	gohttp "net/http"
)

type ResultsVars struct {
	Items []*gofeed.Item
	Query string
}

func SearchHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	// please make this better... (20180407/thisisaaronland)

	q_tpl, err := html.Asset("templates/html/search_query.html")

	if err != nil {
		return nil, err
	}

	qt, err := template.New("search_query").Parse(string(q_tpl))

	if err != nil {
		return nil, err
	}

	r_tpl, err := html.Asset("templates/html/search_results.html")

	if err != nil {
		return nil, err
	}

	rt, err := template.New("search_results").Parse(string(r_tpl))

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

			err = qt.ExecuteTemplate(rsp, "search_query", "")

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

		err = rt.ExecuteTemplate(rsp, "search_results", vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
