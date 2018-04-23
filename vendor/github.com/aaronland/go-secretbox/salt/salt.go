package salt

// this should probably just be put in a generic random string package
// (20180115/thisisaaronland)

import (
	"errors"
	"fmt"
	_ "log"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const min_length int = 16

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

type Salt struct {
	salt string
}

func (s *Salt) String() string {
	return s.salt
}

func (s *Salt) Bytes() []byte {
	return []byte(s.salt)
}

type SaltOptions struct {
	Length int
}

func DefaultSaltOptions() *SaltOptions {

	opts := SaltOptions{
		Length: min_length,
	}

	return &opts
}

func NewRandomSalt(opts *SaltOptions) (*Salt, error) {

	count := len(runes)

	result := make([]string, 0)

	var last string

	for len(result) < opts.Length {

		j := r.Intn(count)
		r := runes[j]
		c := fmt.Sprintf("%c", r)

		if c == last {
			continue
		}

		result = append(result, c)
		last = c
	}

	s := strings.Join(result, "")
	return NewSaltFromString(s)
}

func NewSaltFromString(s string) (*Salt, error) {

	_, err := IsValidSalt(s)

	salt := Salt{s}

	if err != nil {
		return nil, err
	}

	return &salt, nil
}

func IsValidSalt(s string) (bool, error) {

	if len(s) < min_length {
		return false, errors.New("salt is too short")
	}

	return true, nil
}
