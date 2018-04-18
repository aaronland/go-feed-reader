package login

import (
       "errors"
	"github.com/aaronland/go-feed-reader/user"
       "net/http"
)

type Provider interface {
     user.UserDB
     SigninURL() string
     CookieSecret() string
}

func IsLoggedIn(pr Provider, req *http.Request) bool {

     _, err := GetLoggedIn(pr, req)

     if err != nil {
     	return false
     }
     
     return true
}

func GetLoggedIn(pr Provider, req *http.Request) (user.User, error) {

     // cookies := req.Cookies
     
     return nil, errors.New("Please write me")
}
