package user

import (
	"errors"
	"github.com/aaronland/go-feed-reader/password"
	"net/mail"
)

type UserDB interface {
	GetUserById(string) (User, error)
	GetUserByEmail(string) (User, error)
	GetUserByUsername(string) (User, error)
	AddUser(User) error
	// DeleteUser(User) error
	// UpdateUser(User, ...) (User, error)
}

type User interface {
	Id() string
	Username() string
	Email() mail.Address
	Password() password.Password
}

func NewDefaultUser() (User, error) {
	return nil, errors.New("Please write me")
}

type DefaultUser struct {
	User
	id       string
	username string
	email    mail.Address
	password password.Password
}

func (u *DefaultUser) Id() string {
	return u.id
}

func (u *DefaultUser) Username() string {
	return u.username
}

func (u *DefaultUser) Email() mail.Address {
	return u.email
}

func (u *DefaultUser) Password() password.Password {
	return u.password
}
