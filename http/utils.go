package http

import (
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/aaronland/go-feed-reader/crumb"
	"github.com/whosonfirst/go-sanitize"	
	gohttp "net/http"
)

func EnsureLoggedIn(pr login.Provider, rsp gohttp.ResponseWriter, req *gohttp.Request) user.User {

	u, err := login.GetLoggedIn(pr, req)

	if user.IsNotExist(err) {
		gohttp.Redirect(rsp, req, pr.SigninURL(), 303)
		return nil
	}

	if err != nil {
		gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
		return nil
	}

	return u
}

func ValidateCrumb(pr login.Provider, req *gohttp.Request) (string, error) {

	raw_crumb := req.PostFormValue("crumb")

	sn_opts := sanitize.DefaultOptions()

	crumb_var, err := sanitize.SanitizeString(raw_crumb, sn_opts)

	if err != nil {
		return "", err
	}

	if crumb_var == "" {
		return "", err
	}

	ok, err := crumb.ValidateCrumb(crumb_var)

	if err != nil {
		return "", err
	}

	if !ok {
		return "", err
	}

	crumb_var, err = crumb.GenerateCrumb(req)

	if err != nil {
		return "", err
	}

	return crumb_var, nil
}
