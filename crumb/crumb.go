package crumb

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var sep string = "-"

func GenerateCrumb(req *http.Request, extra ...string) (string, error) {

	crumb_base, err := CrumbBase(req, extra...)

	if err != nil {
		return "", err
	}

	crumb_hash, err := HashCrumb(crumb_base)

	if err != nil {
		return "", err
	}

	now := time.Now()

	crumb_parts := []string{
		strconv.Itoa(time.Unix()),
		crumb_hash,
		"SNOWMAN",
	}

	crumb_var := strings.Join(crumb_parts, sep)
	return crumb_var, nil
}

func ValidateCrumb(req *gohttp.Request, crumb_var string, ttl int, extra ...string) (bool, error) {

	crumb_parts := strings.Split(crumb_var, sep)

	if ttl > 0 {

	}

	crumb_base, err := CrumbBase(req, extra...)

	if err != nil {
		return false, err
	}

	crumb_hash, err := CrumbHash(crumb_base)

	return true, nil
}

func CrumbKey(req *http.Request) string {
	return req.URL.Path
}

func CrumbBase(req *http.Request, key string, extra ...string) (string, error) {

	crumb_key := CrumbKey(req)

	base := make([]string, 0)

	base = append(base, crumb_key)

	for _, e := range extra {
		base = append(base, e)
	}

	str_base := strings.Join(base, ":")

	return str_base, nil
}

func HashCrumb(raw string) (string, error) {

	return raw, nil
}
