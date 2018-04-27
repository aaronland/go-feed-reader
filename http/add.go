package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/aaronland/go-feed-reader/crumb"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/arschles/go-bindata-html-template"
	"github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	_ "log"
	gohttp "net/http"
	"net/url"
)

type FormVars struct {
	PageTitle string
	Crumb     string
	User      user.User
}

type PostVars struct {
	PageTitle string
	Crumb     string
	Items     []*gofeed.Item
	User      user.User
}

func AddHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	form_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_feed_form.html",
		"templates/html/inc_foot.html",
	}

	post_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_feed_form.html",
		"templates/html/inc_items.html",
		"templates/html/inc_foot.html",
	}

	t_form, err := template.New("add_form", html.Asset).ParseFiles(form_files...)

	if err != nil {
		return nil, err
	}

	funcs := template.FuncMap{
		"strip_tags": strip.StripTags,
	}

	p_form, err := template.New("add_form", html.Asset).Funcs(funcs).ParseFiles(post_files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		u := EnsureLoggedIn(fr, rsp, req)

		if u == nil {
			return
		}

		switch req.Method {
		case "GET":

			crumb_var, err := crumb.GenerateCrumb(req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars := FormVars{
				PageTitle: "",
				Crumb:     crumb_var,
				User:      u,
			}

			err = t_form.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return

		case "POST":

			feed_url, err := PostString(req, "feed")

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			if feed_url == "" {
				gohttp.Error(rsp, "Missing feed", gohttp.StatusBadRequest)
				return
			}

			_, err = url.Parse(feed_url)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			err = ValidateCrumb(fr, req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			// feed, err := fr.AddFeedForUser(user, feed_url)

			feed, err := fr.AddFeed(feed_url)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			crumb_var, err := crumb.GenerateCrumb(req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars := PostVars{
				PageTitle: "",
				Crumb:     crumb_var,
				Items:     feed.Items,
				User:      u,
			}

			err = p_form.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return

		default:
			gohttp.Error(rsp, "Unsupported method", gohttp.StatusMethodNotAllowed)
			return
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
