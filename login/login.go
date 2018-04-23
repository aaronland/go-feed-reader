package login

import (
       "errors"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/aaronland/go-secretbox"	
       "net/http"
)

type Provider interface {
     user.UserDB
     SigninURL() string
     CookieConfig() CookieConfig     		    
}

type CookieConfig interface {
     Salt() string
     Secret() string
     Name() string     	     
}

func IsLoggedIn(pr Provider, req *http.Request) bool {

     _, err := GetLoggedIn(pr, req)

     if err != nil {
     	return false
     }
     
     return true
}

func GetLoggedIn(pr Provider, req *http.Request) (user.User, error) {

     cfg := pr.CookieConfig()
     
     cookie, err:= req.Cookie(cfg.Name)

     if err != nil {
     	return nil, err
     }
     
	opts := secretbox.NewSecretboxOptions()
	opts.Salt = cfg.Salt()

	sb, err := secretbox.NewSecretbox(cfg.Secret(), opts)

     if err != nil {
     	return nil, err
     }

	body, err := sb.Unlock([]byte(cookie))

     if err != nil {
     	return nil, err
     }

     str_body := string(body)

     parts := strings.Split(str_body, ":")

     user_id := parts[0]
     user_pswd := parts[1]

     u, err := pr.GetUserById(user_id)

     if err != nil {
     	return nil, err
     }

     p := u.Password()

     err = p.Compare(user_pswd)

     if err != nil {
     	return nil, err
     }

     return u, nil
}

