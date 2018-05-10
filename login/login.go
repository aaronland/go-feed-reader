package login

import (
	"errors"
	"fmt"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/aaronland/go-secretbox"
	"github.com/aaronland/go-string/random"
	_ "log"
	"net/http"
	"strings"
)

type Provider interface {
	user.UserDB
	Config() Config
}

type Config interface {
	Cookie() CookieConfig
	URL() URLConfig
}

type URLConfig interface {
	SigninURL() string
	SignupURL() string
	SignoutURL() string
}

type CookieConfig interface {
	Salt() string
	Secret() string
	Name() string
}

func NewDefaultURLConfig() (URLConfig, error) {

	cfg := DefaultURLConfig{
		signin:  "/signin",
		signout: "/signout",
		signup:  "/signup",
	}

	return &cfg, nil
}

type DefaultURLConfig struct {
	URLConfig
	signin  string
	signout string
	signup  string
}

func (c *DefaultURLConfig) SigninURL() string {
	return c.signin
}

func (c *DefaultURLConfig) SignoutURL() string {
	return c.signout
}

func (c *DefaultURLConfig) SignupURL() string {
	return c.signup
}

type DefaultCookieConfig struct {
	CookieConfig
	name   string
	salt   string
	secret string
}

func (c *DefaultCookieConfig) Name() string {
	return c.name
}

func (c *DefaultCookieConfig) Salt() string {
	return c.salt
}

func (c *DefaultCookieConfig) Secret() string {
	return c.secret
}

func NewDefaultCookieConfig() (CookieConfig, error) {

	rand_opts := random.DefaultOptions()
	rand_opts.ASCII = true

	var s string

	s, _ = random.String(rand_opts)
	name := s

	s, _ = random.String(rand_opts)
	salt := s

	s, _ = random.String(rand_opts)
	secret := s

	cfg := DefaultCookieConfig{
		name:   name,
		salt:   salt,
		secret: secret,
	}

	return &cfg, nil
}

type DefaultConfig struct {
	Config
	cookie CookieConfig
	url    URLConfig
}

func (c *DefaultConfig) Cookie() CookieConfig {
	return c.cookie
}

func (c *DefaultConfig) URL() URLConfig {
	return c.url
}

func NewDefaultConfig() (Config, error) {

	c, err := NewDefaultCookieConfig()

	if err != nil {
		return nil, err
	}

	u, err := NewDefaultURLConfig()

	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig{
		cookie: c,
		url:    u,
	}

	return &cfg, nil
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

		if err == http.ErrNoCookie {
			return nil, &user.ErrNoUser{}
		}

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

	if p.Digest() != user_pswd {
		return nil, errors.New("Invalid password")
	}

	return u, nil
}

func GetLoginCookie(pr Provider, req *http.Request) (string, error) {

	cfg := pr.Config()
	ck_cfg := cfg.Cookie()

	cookie, err := req.Cookie(ck_cfg.Name())

	if err != nil {
		return "", err
	}

	opts := secretbox.NewSecretboxOptions()
	opts.Salt = ck_cfg.Salt()

	sb, err := secretbox.NewSecretbox(ck_cfg.Secret(), opts)

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

	cfg := pr.Config()
	ck_cfg := cfg.Cookie()

	opts := secretbox.NewSecretboxOptions()
	opts.Salt = ck_cfg.Salt()

	sb, err := secretbox.NewSecretbox(ck_cfg.Secret(), opts)

	if err != nil {
		return err
	}

	enc, err := sb.Lock([]byte(body))

	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:  ck_cfg.Name(),
		Value: string(enc),
	}

	http.SetCookie(rsp, &cookie)
	return nil
}

func DeleteLoginCookie(pr Provider, rsp http.ResponseWriter) error {

	cfg := pr.Config()
	ck_cfg := cfg.Cookie()

	cookie := http.Cookie{
		Name:   ck_cfg.Name(),
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(rsp, &cookie)
	return nil
}
