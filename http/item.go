package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/arschles/go-bindata-html-template"
	"github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	"github.com/whosonfirst/go-sanitize"
	_ "log"
	gohttp "net/http"
)

type ItemVars struct {
	PageTitle string
	Item      *gofeed.Item
}

func ItemHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_item.html",
		"templates/html/inc_foot.html",
	}

	funcs := template.FuncMap{
		"strip_tags": strip.StripTags,
	}

	t, err := template.New("item", html.Asset).Funcs(funcs).ParseFiles(files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		query := req.URL.Query()

		raw_guid := query.Get("guid")

		sn_opts := sanitize.DefaultOptions()

		str_guid, err := sanitize.SanitizeString(raw_guid, sn_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		/*
		raw_feed := query.Get("feed")
		
		str_feed, err := sanitize.SanitizeString(raw_feed, sn_opts)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}
		*/

		item, err := fr.GetItemByGUID(str_guid)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)		
			return
		}

		vars := ItemVars{
			PageTitle: item.Title,
			Item:      item,
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
