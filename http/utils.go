package http

import (
	"github.com/aaronland/go-feed-reader/crumb"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/user"
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

func ValidateCrumb(pr login.Provider, req *gohttp.Request) error {

	crumb_var, err := PostString(req, "crumb")

	if err != nil {
		return err
	}

	if crumb_var == "" {
		return err
	}

	ok, err := crumb.ValidateCrumb(crumb_var)

	if err != nil {
		return err
	}

	if !ok {
		return err
	}

	return nil
}
