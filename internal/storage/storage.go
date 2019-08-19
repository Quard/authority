package storage

import (
	"errors"

	"github.com/Quard/authority/internal/user"
)

type Storage interface {
	AddUser(user user.User) error
	GetUserByID(ID string) (user.User, error)
	GetUserByEmail(email string) (user.User, error)
	GetUserByProp(name, value string) (user.User, error)
	GetUserBySession(authToken string) (user.User, error)

	AddSession(user user.User, session user.Session) error
	SetUserProp(user user.User, name, value string) error
}

var ErrUserNotFound = errors.New("user not found")
