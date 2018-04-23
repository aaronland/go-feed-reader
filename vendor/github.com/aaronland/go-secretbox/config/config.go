package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func DefaultPathForUser(usr *user.User) string {

	home := usr.HomeDir
	config := filepath.Join(home, ".config")
	root := filepath.Join(config, "secretbox")
	path := filepath.Join(root, "salt")

	return path
}

func ReadSalt(path string) (string, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return "", nil
	}

	_, err = os.Stat(abs_path)

	if err != nil {
		return "", err
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
