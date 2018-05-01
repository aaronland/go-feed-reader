package http

import (
	"github.com/aaronland/go-feed-reader/assets/html"
	"github.com/aaronland/go-feed-reader/crumb"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/user"
	"github.com/arschles/go-bindata-html-template"
	"github.com/grokify/html-strip-tags-go"
	gohttp "net/http"
)

func CompileTemplate(name string, files ...string) (*template.Template, error) {

	funcs := template.FuncMap{
		"strip_tags": strip.StripTags,
	}

	return template.New(name, html.Asset).Funcs(funcs).ParseFiles(files...)
}

func EnsureLoggedIn(pr login.Provider, rsp gohttp.ResponseWriter, req *gohttp.Request) user.User {

	u, err := login.GetLoggedIn(pr, req)

	if user.IsNotExist(err) {
		gohttp.Redirect(rsp, req, pr.SigninURL(), 303)
		return nil
	}

	if err != nil {

		if err == gohttp.ErrNoCookie {
			gohttp.Redirect(rsp, req, pr.SigninURL(), 303)
			return nil
		}

		gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
		return nil
	}

	return u
}

func GenerateCrumb(pr login.Provider, req *gohttp.Request, extra ...string) (string, error) {

	cfg := crumb.NewCrumbConfig()
	return crumb.GenerateCrumb(cfg, req, extra...)
}

func ValidateCrumb(pr login.Provider, req *gohttp.Request) error {

	crumb_var, err := PostString(req, "crumb")

	if err != nil {
		return err
	}

	if crumb_var == "" {
		return err
	}

	cfg := crumb.NewCrumbConfig()

	ok, err := crumb.ValidateCrumb(cfg, req, crumb_var, 0)

	if err != nil {
		return err
	}

	if !ok {
		return err
	}

	return nil
}
