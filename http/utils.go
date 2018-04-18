package http

import (
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/user"
	gohttp "net/http"
)

func EnsureLoggedIn(cfg login.Config, udb user.UserDB, rsp gohttp.ResponseWriter, req *gohttp.Request) (user.User, error) {

	u, err := login.GetLoggedIn(cfg, udb, req)

	if user.IsNotExist(err) {
		rsp.Redirect(cfg.SigninURL())
		return nil
	}

	if err != nil {
		gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
		return nil
	}

	return u
}
