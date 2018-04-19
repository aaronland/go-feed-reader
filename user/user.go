package user

import (
	"github.com/aaronland/go-feed-reader/password"
	_ "net/mail"
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
	Email() string // mail.Address
	Password() password.Password
}

func IsNotExist(e error) bool {
	return false
}

func NewDefaultUser(username string, email string, password password.Password) (User, error) {

	return NewDefaultUserWithID("xxx", username, email, password)
}

func NewDefaultUserWithID(id string, username string, email string, password password.Password) (User, error) {

	u := DefaultUser{
		id:       id,
		username: username,
		email:    email,
		password: password,
	}

	return &u, nil
}

type DefaultUser struct {
	User
	id       string
	username string
	email    string // mail.Address
	password password.Password
}

func (u *DefaultUser) Id() string {
	return u.id
}

func (u *DefaultUser) Username() string {
	return u.username
}

func (u *DefaultUser) Email() string { // mail.Address {
	return u.email
}

func (u *DefaultUser) Password() password.Password {
	return u.password
}
