package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/aaronland/go-feed-reader/crumb"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/arschles/go-bindata-html-template"
	"github.com/whosonfirst/go-sanitize"
	_ "log"
	gohttp "net/http"
	// "net/mail"
)

type SigninVars struct {
	PageTitle string
	Crumb     string
	Error     error
}

func SigninHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	form_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_signin_form.html",
		"templates/html/inc_foot.html",
	}

	s_form, err := template.New("signin_form", html.Asset).ParseFiles(form_files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		if login.IsLoggedIn(fr, req) {
			gohttp.Redirect(rsp, req, "/", 303)
			return
		}

		switch req.Method {
		case "GET":

			crumb_var, err := crumb.GenerateCrumb(req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars := SigninVars{
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

			raw_email := req.PostFormValue("email")
			raw_password := req.PostFormValue("password")
			raw_crumb := req.PostFormValue("crumb")

			opts := sanitize.DefaultOptions()

			str_email, err := sanitize.SanitizeString(raw_email, opts)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			str_password, err := sanitize.SanitizeString(raw_password, opts)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			/* sudo put me in a function */

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

			crumb_var, err = crumb.GenerateCrumb(req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			/* end of sudo put me in a function */

			u, err := fr.GetUserByEmail(str_email)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			p := u.Password()

			err = p.Compare(str_password)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			gohttp.Redirect(rsp, req, "/", 303)
			return

		default:
			gohttp.Error(rsp, "Unsupported method", gohttp.StatusMethodNotAllowed)
			return
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
