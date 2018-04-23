package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-secretbox/config"
	"github.com/aaronland/go-secretbox/salt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func main() {

	var output = flag.String("output", "", "Where to write the salt file. Default is ~/.config/secretbox/salt")
	var length = flag.Int("length", 16, "How long the salt should be")

	flag.Parse()

	if *output == "" {

		usr, err := user.Current()

		if err != nil {
			log.Fatal(err)
		}

		*output = config.DefaultPathForUser(usr)

		root := filepath.Dir(*output)
		_, err = os.Stat(root)

		if os.IsNotExist(err) {

			err = os.MkdirAll(root, 0700)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	opts := salt.DefaultSaltOptions()
	opts.Length = *length

	s, err := salt.NewRandomSalt(opts)

	if err != nil {
		log.Fatal(err)
	}

	var writer io.Writer

	if *output == "STDOUT" {
		writer = os.Stdout
	} else {

		info, err := os.Stat(*output)

		if err == nil {

			mtime := info.ModTime()
			ts := mtime.Unix()

			root := filepath.Dir(*output)
			fname := filepath.Base(*output)

			fname = fmt.Sprintf("%s.%d", fname, ts)
			backup := filepath.Join(root, fname)

			err = os.Rename(*output, backup)

			if err != nil {
				log.Fatal(err)
			}
		}

		fh, err := os.OpenFile(*output, os.O_RDWR|os.O_CREATE, 0600)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()
		writer = fh
	}

	_, err = writer.Write(s.Bytes())

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
