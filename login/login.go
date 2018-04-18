package login

import (
       "errors"
	"github.com/aaronland/go-feed-reader/user"
       "net/http"
)

func IsLoggedIn(udb user.UserDB, req *http.Request) bool {

     _, err := GetLoggedIn(udb, req)

     if err != nil {
     	return false
     }
     
     return true
}

func GetLoggedIn(udb user.UserDB, req *http.Request) (user.User, error) {

     // cookies := req.Cookies
     
     return nil, errors.New("Please write me")
}
