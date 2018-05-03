package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-string/random"
	"log"
)

func main() {

	var ascii = flag.Bool("ascii", false, "")
	var length = flag.Int("length", 32, "")
	var chars = flag.Int("chars", 0, "")

	flag.Parse()

	opts := random.DefaultOptions()
	opts.ASCII = *ascii
	opts.Length = *length
	opts.Chars = *chars

	s, err := random.String(opts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
}
