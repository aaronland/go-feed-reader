package user

import (
       "errors"
       "github.com/aaronland/go-feed-reader/password"	
       "net/mail"
)

type User interface {
     Username() string
     Email() mail.Address     	      
     Password() password.Password
}

func NewDefaultUser() (User, error) {
     return errors.New("Please write me")
}

type DefaultUser struct {
     User
     username string
     email mail.Address
     password password.Password
}

func (u *DefaultUser) Username() {
     return u.username
}

func (u *DefaultUser) Email() {
     return u.email
}

func (u *DefaultUser) Password() {
     return u.password
}
