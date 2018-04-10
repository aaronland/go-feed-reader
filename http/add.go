package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/arschles/go-bindata-html-template"
	"github.com/whosonfirst/go-sanitize"
	"github.com/grokify/html-strip-tags-go"	
	"log"		
	gohttp "net/http"
	"net/url"
)

type FormVars struct {
     	PageTitle string
	Crumb string
}

func AddHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	form_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_feed_form.html",
		"templates/html/inc_foot.html",
	}

	post_files := []string{
		"templates/html/inc_head.html",
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

		log.Println("ADD", req.Method)
		
		switch req.Method {
		case "GET":

			vars := FormVars{
			     	PageTitle: "",
				Crumb: "OMGWTFFIXME",
			}

			err := t_form.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return						  
			}
			
			return
			
		case "POST":

			log.Println("FORM")
			
			raw_feed := req.PostFormValue("feed")
			raw_crumb := req.PostFormValue("crumb")			

			opts := sanitize.DefaultOptions()
			
			feed_url, err := sanitize.SanitizeString(raw_reed, opts)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return						  
			}

			crumb, err := sanitize.SanitizeString(raw_crumb, opts)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return						  
			}

			if feed_url == "" {
				gohttp.Error(rsp, "Missing feed", gohttp.StatusBadRequest)
				return						  
			}

			if crumb == "" {
				gohttp.Error(rsp, "Missing crumb", gohttp.StatusBadRequest)
				return						  
			}

			_, err = url.Parse(feed_url)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return						  
			}

			_, err = fr.AddFeed(feed_url)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return						  
			}

			// display stuff here...
			return

		default:
			gohttp.Error(rsp, "Unsupported method", gohttp.StatusMethodNotAllowed)
			return
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
