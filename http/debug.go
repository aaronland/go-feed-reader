package http

import (
	"github.com/aaronland/go-feed-reader"
	_ "log"
	gohttp "net/http"
)

func DebugHandler(fr *reader.FeedReader) (gohttp.Handler, error) {

	files := []string{
		"templates/atom/debug.xml",
	}

	_, err := CompileTemplate("feed", files...)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
