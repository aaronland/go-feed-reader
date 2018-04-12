package http

import (
	"github.com/whosonfirst/go-sanitize"
	gohttp "net/http"
)

func GetString(req *gohttp.Request, param string) (string, error) {

	opts := sanitize.DefaultOptions()

	q := req.URL.Query()
	raw_p := q.Get(param)

	return sanitize.SanitizeString(raw_p, opts)
}
