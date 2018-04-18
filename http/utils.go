package http

import (
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/user"
	gohttp "net/http"
)

func EnsureLoggedIn(udb user.UserDB, rsp gohttp.ResponseWriter, req *gohttp.Request) (user.User, error) {

	u, err := login.GetLoggedIn(udb, req)

	if user.IsNotExist(err) {
		rsp.Redirect("/signin")
		return nil
	}

	if err != nil {
		gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
		return nil
	}

	return u
}
