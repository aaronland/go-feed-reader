package login

import (
	_ "errors"
	"fmt"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/aaronland/go-secretbox"
	"net/http"
	"strings"
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

	cookie, err := GetLoginCookie(pr, req)

	if err != nil {
		return nil, err
	}

	parts := strings.Split(cookie, ":")

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

func GetLoginCookie(pr Provider, req *http.Request) (string, error) {

	cfg := pr.CookieConfig()

	cookie, err := req.Cookie(cfg.Name())

	if err != nil {
		return "", err
	}

	opts := secretbox.NewSecretboxOptions()
	opts.Salt = cfg.Salt()

	sb, err := secretbox.NewSecretbox(cfg.Secret(), opts)

	if err != nil {
		return "", err
	}

	body, err := sb.Unlock([]byte(cookie.Value))

	if err != nil {
		return "", err
	}

	str_body := string(body)

	return str_body, nil
}

func SetLoginCookie(pr Provider, rsp http.ResponseWriter, u user.User) error {

	pswd := u.Password()
	body := fmt.Sprintf("%s:%s", u.Id(), pswd.Digest())

	cfg := pr.CookieConfig()

	opts := secretbox.NewSecretboxOptions()
	opts.Salt = cfg.Salt()

	sb, err := secretbox.NewSecretbox(cfg.Secret(), opts)

	if err != nil {
		return err
	}

	enc, err := sb.Lock([]byte(body))

	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:  cfg.Name(),
		Value: string(enc),
	}

	http.SetCookie(rsp, &cookie)
	return nil
}
