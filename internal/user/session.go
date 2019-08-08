package user

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Session struct {
	AuthToken   string
	CreatedDate time.Time
}

func (u User) CreateSession() (Session, error) {
	hash := sha256.New()
	hash.Write([]byte(time.Now().String()))
	hash.Write([]byte(u.Email))

	session := Session{
		AuthToken:   fmt.Sprintf("%x", hash.Sum(nil)),
		CreatedDate: time.Now(),
	}
	u.Sessions = append(u.Sessions, session)

	return session, nil
}
