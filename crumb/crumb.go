package crumb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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
		Secret:  "secretsdfsjdfhskdhfskdhfskjdfhw4",
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

	crumb_var := fmt.Sprintf("%s-%s", str_ts, crumb_hash)

	enc_var, err := Encrypt(cfg, crumb_var)

	if err != nil {
		return "", err
	}

	return enc_var, nil
}

func ValidateCrumb(cfg CrumbConfig, req *http.Request, enc_var string, ttl int64, extra ...string) (bool, error) {

	crumb_var, err := Decrypt(cfg, enc_var)

	if err != nil {
		return false, err
	}

	crumb_parts := strings.Split(crumb_var, "-")

	if len(crumb_parts) != 2 {
		return false, errors.New("Invalid crumb")
	}

	crumb_ts := crumb_parts[0]
	crumb_test := crumb_parts[1]

	crumb_base, err := CrumbBase(cfg, req, extra...)

	if err != nil {
		return false, err
	}

	crumb_hash, err := HashCrumb(cfg, crumb_base)

	if err != nil {
		return false, err
	}

	ok, err := CompareHashes(crumb_hash, crumb_test)

	if err != nil {
		return false, err
	}

	if !ok {
		return false, errors.New("Crumb mismatch")
	}

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

	return true, nil
}

func CrumbKey(cfg CrumbConfig, req *http.Request) string {
	return fmt.Sprintf("%s-%s", cfg.Snowman, req.URL.Path)
}

func CrumbBase(cfg CrumbConfig, req *http.Request, extra ...string) (string, error) {

	crumb_key := CrumbKey(cfg, req)

	base := make([]string, 0)

	base = append(base, crumb_key)
	base = append(base, req.UserAgent())
	// base = append(base, req.RemoteAddr)

	for _, e := range extra {
		base = append(base, e)
	}

	str_base := strings.Join(base, "-")
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

// https://gist.github.com/manishtpatel/8222606
// https://github.com/blaskovicz/go-cryptkeeper/blob/master/encrypted_string.go

func Encrypt(cfg CrumbConfig, text string) (string, error) {

	plaintext := []byte(text)
	secret := []byte(cfg.Secret)

	block, err := aes.NewCipher(secret)

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(cfg CrumbConfig, cryptoText string) (string, error) {

	ciphertext, err := hex.DecodeString(cryptoText)

	if err != nil {
		return "", err
	}

	secret := []byte(cfg.Secret)
	block, err := aes.NewCipher(secret)

	if err != nil {
		return "", err
	}

	if byteLen := len(ciphertext); byteLen < aes.BlockSize {
		return "", fmt.Errorf("invalid cipher size %d.", byteLen)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
