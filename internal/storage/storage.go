package storage

import (
	"errors"

	"github.com/Quard/authority/internal/user"
)

type Storage interface {
	AddUser(user user.User) error
	GetUserByEmail(email string) (user.User, error)
}

var ErrUserNotFound = errors.New("user not found")
