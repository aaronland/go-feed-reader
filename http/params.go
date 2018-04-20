package http

import (
	"github.com/whosonfirst/go-sanitize"
	gohttp "net/http"
	"strconv"
)

var sn_opts *sanitize.Options

func init() {

	sn_opts = sanitize.DefaultOptions()
}

func GetString(req *gohttp.Request, param string) (string, error) {

     q := req.URL.Query()
     raw_value := q.Get(param)
     return sanitize.SanitizeString(raw_value, sn_opts)		
}

func PostString(req *gohttp.Request, param string) (string, error) {

     raw_value := req.PostForm(param)
     return sanitize.SanitizeString(raw_value, sn_opts)		
}

func GetInt64(req *gohttp.Request, param string) (int64, error){

     str_value, err := GetString(req, param)

     if err != nil {
     	return -1, err
     }

     return strconv.ParseInt(str_value, 10, 64)
}

func PostInt64(req *gohttp.Request, param string) (int64, error){

     str_value, err := PostString(req, param)

     if err != nil {
     	return -1, err
     }

     return strconv.ParseInt(str_value, 10, 64)
}
