package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/aaronland/go-feed-reader/crumb"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/arschles/go-bindata-html-template"
	_ "log"
	gohttp "net/http"
)

type SignoutVars struct {
	PageTitle string
	Crumb     string
	Error     error
}

func SignoutHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	form_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_signout_form.html",
		"templates/html/inc_foot.html",
	}

	s_form, err := template.New("signin_form", html.Asset).ParseFiles(form_files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		user := EnsureLoggedIn(fr, rsp, req)

		if user == nil {
		   	return 
		}

		switch req.Method {
		case "GET":

			crumb_var, err := crumb.GenerateCrumb(req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars := SignoutVars{
				PageTitle: "",
				Crumb:     crumb_var,
				Error:     nil,
			}

			err = s_form.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return

		case "POST":

			err := ValidateCrumb(fr, req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			err = login.DeleteLoginCookie(fr, rsp)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			gohttp.Redirect(rsp, req, "/")
			return

		default:
			gohttp.Error(rsp, "Unsupported method", gohttp.StatusMethodNotAllowed)
			return
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
