package password

import ()

type Password interface {
	Digest() string
	Compare(string) error
}
