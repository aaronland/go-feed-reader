package main

// https://godoc.org/golang.org/x/crypto/nacl/secretbox

// please don't import anything that isn't part of the standard
// library or "golang.org/x/" unless there's a really good reason
// to (20171025/thisisaaronland)

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-secretbox"
	"github.com/aaronland/go-secretbox/config"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {

	var suffix = flag.String("suffix", ".enc", "...")
	var unlock = flag.Bool("unlock", false, "Decrypt files.")
	var debug = flag.Bool("debug", false, "...")
	var salt = flag.String("salt", "config:", "...")

	flag.Parse()

	possible := flag.Args()
	files := make([]string, 0)

	if len(possible) == 0 {
		log.Println("No secrets to tell!")
		os.Exit(0)
	}

	for _, path := range possible {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		files = append(files, abs_path)
	}

	if *salt == "env:" {
		*salt = os.Getenv("SECRETBOX_SALT")
	} else if strings.HasPrefix(*salt, "config:") {

		parts := strings.Split(*salt, ":")
		var path string

		if parts[1] == "" {
			usr, err := user.Current()

			if err != nil {
				log.Fatal(err)
			}

			path = config.DefaultPathForUser(usr)

		} else {
			path = parts[1]
		}

		s, err := config.ReadSalt(path)

		if err != nil {
			log.Fatal(err)
		}

		*salt = s

	} else if *salt != "" {
		// pass
	} else {
		log.Fatal("missing salt")
	}

	if len(*salt) < 8 {
		log.Fatal("invalid salt")
	}

	fmt.Println("enter password: ")
	pswd, err := terminal.ReadPassword(0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("enter password (again): ")
	pswd2, err := terminal.ReadPassword(0)

	if err != nil {
		log.Fatal(err)
	}

	if string(pswd) != string(pswd2) {
		log.Fatal("password mismatch")
	}

	opts := secretbox.NewSecretboxOptions()
	opts.Salt = *salt
	opts.Suffix = *suffix
	opts.Debug = *debug

	sb, err := secretbox.NewSecretbox(string(pswd), opts)

	for _, abs_path := range files {

		var sb_path string
		var sb_err error

		if *unlock {
			sb_path, sb_err = sb.UnlockFile(abs_path)
		} else {
			sb_path, sb_err = sb.LockFile(abs_path)
		}

		if sb_err != nil {
			log.Fatal(sb_err)
		}

		log.Println(sb_path)
	}

	os.Exit(0)
}
