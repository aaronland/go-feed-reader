package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/mmcdole/gofeed"
	_ "log"
	gohttp "net/http"
)

type ItemVars struct {
	PageTitle string
	Item      *gofeed.Item
	Feed      *gofeed.Feed
	User      user.User
}

func ItemHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_item.html",
		"templates/html/inc_foot.html",
	}

	t, err := CompileTemplate("item", files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		u := EnsureLoggedIn(fr, rsp, req)

		if u == nil {
			return
		}

		str_guid, err := GetString(req, "guid")

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

		f, err := fr.GetFeedByItemGUID(str_guid)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		vars := ItemVars{
			PageTitle: item.Title,
			Item:      item,
			Feed:      f,
			User:      u,
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
