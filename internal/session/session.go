package session

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Quard/authority/internal/user"
)

type Session struct {
	AuthToken string
	User      user.User
}

func CreateSession(user user.User) (Session, error) {
	hash := sha256.New()
	hash.Write([]byte(time.Now().String()))
	hash.Write([]byte(user.Email))

	session := Session{
		AuthToken: fmt.Sprintf("%x", hash.Sum(nil)),
		User:      user,
	}

	return session, nil
}
