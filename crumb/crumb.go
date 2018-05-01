package crumb

import (
	"errors"
	"github.com/aaronland/go-secretbox"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sep string = "-"

type CrumbProvider interface {
	GenerateCrumb(*http.Request, ...string) (string, error)
	ValidateCrumb(*http.Request, string, int64, ...string) (bool, error)
}

type CrumbConfig struct {
	Snowman string
	Salt    string
	Secret  string
}

func NewCrumbConfig() CrumbConfig {

	cfg := CrumbConfig{
		Snowman: "SNOWMAN",
		Salt:    "salt",
		Secret:  "secret",
	}

	return cfg
}

func GenerateCrumb(cfg CrumbConfig, req *http.Request, extra ...string) (string, error) {

	crumb_base, err := CrumbBase(cfg, req, extra...)

	if err != nil {
		return "", err
	}

	crumb_hash, err := HashCrumb(cfg, crumb_base)

	if err != nil {
		return "", err
	}

	now := time.Now()
	// ts := now.Format(time.Unix)

	ts := now.Unix()

	crumb_parts := []string{
		strconv.FormatInt(ts, 10),
		crumb_hash,
		cfg.Snowman,
	}

	crumb_var := strings.Join(crumb_parts, sep)
	return crumb_var, nil
}

func ValidateCrumb(cfg CrumbConfig, req *http.Request, crumb_var string, ttl int64, extra ...string) (bool, error) {

	crumb_parts := strings.Split(crumb_var, sep)

	if ttl > 0 {

		then, err := strconv.ParseInt(crumb_parts[0], 10, 64)

		if err != nil {
			return false, err
		}

		now := time.Now()
		ts := now.Unix()

		if ts-then > ttl {
			return false, errors.New("Crumb has expired")
		}
	}

	crumb_base, err := CrumbBase(cfg, req, extra...)

	if err != nil {
		return false, err
	}

	crumb_hash, err := HashCrumb(cfg, crumb_base)

	if len(crumb_hash) != len(crumb_var) {
		return false, errors.New("Invalid crumb")
	}

	if crumb_hash != crumb_var {
		return false, errors.New("Invalid crumb")
	}

	return true, nil
}

func CrumbKey(cfg CrumbConfig, req *http.Request) string {
	return req.URL.Path
}

func CrumbBase(cfg CrumbConfig, req *http.Request, extra ...string) (string, error) {

	crumb_key := CrumbKey(cfg, req)

	base := make([]string, 0)

	base = append(base, crumb_key)

	for _, e := range extra {
		base = append(base, e)
	}

	str_base := strings.Join(base, ":")

	return str_base, nil
}

func HashCrumb(cfg CrumbConfig, raw string) (string, error) {

	opts := secretbox.NewSecretboxOptions()
	opts.Salt = cfg.Salt

	sb, err := secretbox.NewSecretbox(cfg.Secret, opts)

	enc, err := sb.Lock([]byte(raw))

	if err != nil {
		return "", err
	}

	return enc, nil
}
