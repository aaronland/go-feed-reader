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
	Email() string // mail.Address
	Password() password.Password
}

func IsNotExist(e error) bool {
	return false
}

func NewDefaultUserRaw(udb UserDB, username string, email string, pswd password.Password) (User, error) {

	if len(username) < 3 {
		return nil, errors.New("Username too short")
	}

	if len(username) > 255 {
		return nil, errors.New("Username too long")
	}

	e, err := mail.ParseAddress(email)

	if err != nil {
		return nil, err
	}

	email = e.Address

	u, err := udb.GetUserByEmail(email)

	if err != nil && !IsNotExist(err) {
		return nil, err
	}

	if u != nil {
		return nil, errors.New("user already exists")
	}

	u, err = udb.GetUserByUsername(username)

	if err != nil && !IsNotExist(err) {
		return nil, err
	}

	if u != nil {
		return nil, errors.New("user already exists")
	}

	u, err = NewDefaultUser(username, email, pswd)

	if err != nil {
		return nil, err
	}

	return u, nil
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
