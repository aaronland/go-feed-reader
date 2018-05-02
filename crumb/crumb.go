package crumb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	_ "fmt"
	_ "log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sep string = "-"

type CrumbConfig struct {
	Snowman string
	Secret  string
}

func NewCrumbConfig() CrumbConfig {

	cfg := CrumbConfig{
		Snowman: "SNOWMAN",
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
	ts := now.Unix()

	str_ts := strconv.FormatInt(ts, 10)

	crumb_parts := []string{
		str_ts,
		crumb_hash,
		cfg.Snowman,
	}

	crumb_var := strings.Join(crumb_parts, sep)
	return crumb_var, nil
}

func ValidateCrumb(cfg CrumbConfig, req *http.Request, crumb_var string, ttl int64, extra ...string) (bool, error) {

	crumb_parts := strings.Split(crumb_var, sep)

	if len(crumb_parts) != 3 {
		return false, errors.New("Malformed crumb")
	}

	crumb_ts := crumb_parts[0]
	crumb_hash := crumb_parts[1]

	if ttl > 0 {

		then, err := strconv.ParseInt(crumb_ts, 10, 64)

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

	crumb_test, err := HashCrumb(cfg, crumb_base)

	if err != nil {
		return false, err
	}

	return CompareHashes(crumb_hash, crumb_test)
}

func CrumbKey(cfg CrumbConfig, req *http.Request) string {
	return req.URL.Path
}

func CrumbBase(cfg CrumbConfig, req *http.Request, extra ...string) (string, error) {

	crumb_key := CrumbKey(cfg, req)

	base := make([]string, 0)

	base = append(base, crumb_key)
	base = append(base, req.UserAgent())

	for _, e := range extra {
		base = append(base, e)
	}

	str_base := strings.Join(base, ":")

	return str_base, nil
}

func CompareHashes(this_enc string, that_enc string) (bool, error) {

	this_hash, err := hex.DecodeString(this_enc)

	if err != nil {
		return false, err
	}

	that_hash, err := hex.DecodeString(that_enc)

	if err != nil {
		return false, err
	}

	match := hmac.Equal(this_hash, that_hash)
	return match, nil
}

func HashCrumb(cfg CrumbConfig, raw string) (string, error) {

	msg := []byte(raw)
	key := []byte(cfg.Secret)

	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	hash := mac.Sum(nil)

	enc := hex.EncodeToString(hash[:])
	return enc, nil
}
