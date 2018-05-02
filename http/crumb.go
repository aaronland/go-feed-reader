package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/user"
	_ "log"
	gohttp "net/http"
)

type CrumbVars struct {
	PageTitle string
	Crumb     string
	User      user.User
	Error     error
}

func CrumbHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	crumb_files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_crumb.html",
		"templates/html/inc_foot.html",
	}

	t, err := CompileTemplate("crumb", crumb_files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		switch req.Method {
		case "GET":

			crumb_var, err := GenerateCrumb(fr, req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars := CrumbVars{
				PageTitle: "",
				Crumb:     crumb_var,
				User:      nil,
				Error:     nil,
			}

			err = t.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return

		case "POST":

			v_err := ValidateCrumb(fr, req)

			crumb_var, err := GenerateCrumb(fr, req)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			vars := CrumbVars{
				PageTitle: "",
				Crumb:     crumb_var,
				Error:     v_err,
				User:      nil,
			}

			err = t.Execute(rsp, vars)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			return

		default:
			gohttp.Error(rsp, "Unsupported method", gohttp.StatusMethodNotAllowed)
			return
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
