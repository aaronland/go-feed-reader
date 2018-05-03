package random

import (
	_ "errors"
	"fmt"
	_ "log"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const min_length int = 32

var runes []rune

var r *rand.Rand

func init() {

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	runes = make([]rune, 0)

	codepoints := [][]int{
		[]int{1, 255},         // ascii
		[]int{127744, 128317}, // emoji
	}

	for _, r := range codepoints {

		first := r[0]
		last := r[1]

		for i := first; i < last; i++ {

			r := rune(i)

			if unicode.IsControl(r) {
				continue
			}

			if unicode.IsSpace(r) {
				continue
			}

			if unicode.IsMark(r) {
				continue
			}

			runes = append(runes, r)
		}
	}

}

type Options struct {
	Length int
	Chars  int
	ASCII  bool
}

func DefaultOptions() *Options {

	opts := Options{
		Length: min_length,
		Chars:  0,
		ASCII:  false,
	}

	return &opts
}

func String(opts *Options) (string, error) {

	count := len(runes)

	result := make([]string, 0)

	var last string

	// chars := 0
	b := 0

	for b < opts.Length {

		j := r.Intn(count)
		r := runes[j]

		if opts.ASCII && r > 127 {
			continue
		}

		c := fmt.Sprintf("%c", r)

		if c == last {
			continue
		}

		last = c

		b += len(c)

		if b <= opts.Length {
			result = append(result, c)
		} else {

			if len(result) > 2 {
				result = result[0 : len(result)-2]
			} else {
				result = make([]string, 0)
			}
			b = len(strings.Join(result, ""))
		}
	}

	s := strings.Join(result, "")
	return s, nil
}
