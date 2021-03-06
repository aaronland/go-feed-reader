package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/password"
	"github.com/aaronland/go-feed-reader/user"
	gohttp "net/http"
)

type SignupVars struct {
	PageTitle string
	Crumb     string
	User      user.User
}

func SignupHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	form_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_signup_form.html",
		"templates/html/inc_foot.html",
	}

	s_form, err := CompileTemplate("add_form", form_files...)

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

			vars := SignupVars{
				PageTitle: "",
				Crumb:     crumb_var,
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

			str_username, err := PostString(req, "username")

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

			/* make user stuff */

			salt := "FIXME"
			pswd, err := password.NewBCryptPassword(str_password, salt)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			u, err := user.NewDefaultUserRaw(fr, str_username, str_email, pswd)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			err = fr.AddUser(u)

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
