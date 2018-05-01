package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/user"
	gohttp "net/http"
)

type SigninVars struct {
	PageTitle string
	Crumb     string
	Error     error
	User      user.User
}

func SigninHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	form_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_signin_form.html",
		"templates/html/inc_foot.html",
	}

	s_form, err := CompileTemplate("signin_form", form_files...)

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

			crumb_var, err := GenerateCrumb(fr, req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars := SigninVars{
				PageTitle: "",
				Crumb:     crumb_var,
				Error:     nil,
				User:      nil,
			}

			err = s_form.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return

		case "POST":

			str_email, err := PostString(req, "email")

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			str_password, err := PostString(req, "password")

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
				return
			}

			err = ValidateCrumb(fr, req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

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

			err = login.SetLoginCookie(fr, rsp, u)

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
