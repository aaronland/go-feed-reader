package password

import (
	"crypto/sha512"
	"github.com/patrickmn/go-hmaccrypt"
)

type BCryptPassword struct {
	Password
	crypt *hmaccrypt.HmacCrypt		
	digest string
}

func NewBCryptPassword(pswd string) (Password, error) {

	pepper := []byte("randomly generated sequence stored on disk or in the source")
	crypt := hmaccrypt.New(sha512.New, pepper)

	b_pswd := []byte(pswd)
	digest, err := crypt.Bcrypt(b_pswd, 10)

	if err != nil {
		return nil, err
	}

	p := BCryptPassword{
		digest: digest,
		crypt: crypt,
	}

	return &p, nil
}

func (p *BCryptPassword) Compare(pswd string) error {

	return p.crypt.BcryptCompare(p.digest, pswd)
}
