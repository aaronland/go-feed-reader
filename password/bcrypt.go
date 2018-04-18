package password

import (
	"crypto/sha512"
	"github.com/pmylund/go-hmaccrypt"
)

type BCryptPassword struct {
	Password
	digest string
}

func NewBCryptPasswordFromDigest(digest string) (Password, error) {

	p := BCryptPassword{
		digest: digest,
	}

	return &p, nil
}

func NewBCryptPassword(pswd string) (Password, error) {

	pepper := []byte("randomly generated sequence stored on disk or in the source")
	crypt := hmaccrypt.New(sha512.New, pepper)

	b_pswd := []byte(pswd)

	digest, err := crypt.Bcrypt(b_pswd, 10)

	if err != nil {
		return nil, err
	}

	return NewBCryptPasswordFromDigest(digest)
}

func (p *BCryptPassword) Compare(pswd string) error {

	return crypt.BcryptCompare(p.diget, pswd)
}
