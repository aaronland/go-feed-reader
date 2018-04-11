package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/aaronland/go-feed-reader/crumb"	
	"github.com/arschles/go-bindata-html-template"
	"github.com/whosonfirst/go-sanitize"
	"github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"	
	_ "log"		
	gohttp "net/http"
	"net/url"
)

type FormVars struct {
     	PageTitle string
	Crumb string
}

type PostVars struct {
     	PageTitle string
	Crumb string
	Items []*gofeed.Item
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

		/*
		user, err := login.EnsureLogin(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusForbidden)
			return						  
		}
		*/
		
		switch req.Method {
		case "GET":

			crumb_var, err := crumb.GenerateCrumb(req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return						  
			}

			vars := FormVars{
			     	PageTitle: "",
				Crumb: crumb_var,
			}

			err = t_form.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return						  
			}
			
			return
			
		case "POST":

			raw_feed := req.PostFormValue("feed")
			raw_crumb := req.PostFormValue("crumb")			

			opts := sanitize.DefaultOptions()
			
			feed_url, err := sanitize.SanitizeString(raw_feed, opts)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return						  
			}

			crumb_var, err := sanitize.SanitizeString(raw_crumb, opts)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return						  
			}

			if crumb_var == "" {
				gohttp.Error(rsp, "Missing crumb", gohttp.StatusBadRequest)
				return						  
			}

			ok, err := crumb.ValidateCrumb(crumb_var)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return						  
			}

			if !ok {
				gohttp.Error(rsp, "Invalid crumb", gohttp.StatusForbidden)
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

			// feed, err := fr.AddFeedForUser(user, feed_url)
			
			feed, err := fr.AddFeed(feed_url)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return						  
			}

			crumb_var, err = crumb.GenerateCrumb(req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return						  
			}

			vars := PostVars{
			     	PageTitle: "",
				Crumb: crumb_var,
				Items: feed.Items,
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
