package http

import (
	"github.com/whosonfirst/go-sanitize"
	gohttp "net/http"
	"strconv"
)

func GetString(req *gohttp.Request, param string) (string, error) {

	opts := sanitize.DefaultOptions()

	q := req.URL.Query()
	raw_p := q.Get(param)

	return sanitize.SanitizeString(raw_p, opts)
}

func GetInt(req *gohttp.Request, param string) (int, error) {

	p, err := GetString(req, param)

	if err != nil {
		return -1, err
	}

	return strconv.Atoi(p)
}
