package password

import (
	"crypto/sha512"
	"github.com/patrickmn/go-hmaccrypt"
)

type BCryptPassword struct {
	Password
	crypt  *hmaccrypt.HmacCrypt
	digest string
}

func NewBCryptPasswordFromDigest(digest string, salt string) (Password, error) {

	pepper := []byte(salt)
	crypt := hmaccrypt.New(sha512.New, pepper)

	p := BCryptPassword{
		digest: digest,
		crypt:  crypt,
	}

	return &p, nil
}

func NewBCryptPassword(pswd string, salt string) (Password, error) {

	pepper := []byte(salt)
	crypt := hmaccrypt.New(sha512.New, pepper)

	b_pswd := []byte(pswd)
	digest, err := crypt.Bcrypt(b_pswd, 10)

	if err != nil {
		return nil, err
	}

	p := BCryptPassword{
		digest: string(digest),
		crypt:  crypt,
	}

	return &p, nil
}

func (p *BCryptPassword) Digest() string {
     return p.digest
}

func (p *BCryptPassword) Compare(pswd string) error {

	return p.crypt.BcryptCompare([]byte(p.digest), []byte(pswd))
}
