package crumb

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func GenerateCrumb(req *http.Request, extra ...string) (string, error) {

	key := req.URL.Path

	crumb_base, err := CrumbBase(req, key, extra...)

	if err != nil {
		return "", err
	}

	crumb_hash, err := HashCrumb(crumb_base)

	if err != nil {
		return "", err
	}

	now := time.Now()

	crumb_var := fmt.Sprintf("%d-%s-SNOWMAN", now.Unix(), crumb_hash)

	return crumb_var, nil
}

func ValidateCrumb(crumb_var string) (bool, error) {

	return true, nil
}

func CrumbBase(req *http.Request, key string, extra ...string) (string, error) {

	base := make([]string, 0)

	base = append(base, key)

	for _, e := range extra {
		base = append(base, e)
	}

	str_base := strings.Join(base, ":")

	return str_base, nil
}

func HashCrumb(raw string) (string, error) {

	return raw, nil
}
