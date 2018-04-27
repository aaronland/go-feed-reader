package http

import (
	"github.com/aaronland/go-feed-reader"
	"github.com/aaronland/go-feed-reader/login"
	"github.com/aaronland/go-feed-reader/user"
	"log"
	gohttp "net/http"
)

type IndexVars struct {
	PageTitle string
	User      user.User
}

func IndexHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/html/inc_head.html",
		"templates/html/inc_index.html",
		"templates/html/inc_foot.html",
	}

	t, err := CompileTemplate("index", files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		log.Println("INDEX", req.URL)
		log.Println("HOST ", req.URL.Host)
		log.Println("PATH", req.URL.Path)

		u, err := login.GetLoggedIn(fr, req)

		if err != nil && !user.IsNotExist(err) {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		if u != nil {
			gohttp.Redirect(rsp, req, "/recent", 303)
			return
		}

		vars := IndexVars{
			PageTitle: "",
			User:      u,
		}

		err = t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
