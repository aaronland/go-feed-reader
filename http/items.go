package http

import (
	"fmt"
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"	
	"github.com/mmcdole/gofeed"
	"html/template"
	gohttp "net/http"
)

type HTMLVars struct {
	Items []*gofeed.Item
}

func ItemsHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

     	tpl, err := html.Asset("templates/html/items.html")

	if err != nil {
		return nil, err
	}

	t, err := template.New("feeds").Parse(string(tpl))

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		feeds, err := fr.ListItems()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := HTMLVars{
			Items: items,
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
