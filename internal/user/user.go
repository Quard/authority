package user

import (
	"crypto/rand"
	"crypto/sha512"
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Email      string             `bson:"email"`
	Salt       []byte             `bson:"salt"`
	Password   []byte             `bson:"password"`
	DateJoined time.Time          `json:"date_joined"`
}

func (u *User) SetPassword(password string) error {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		sentry.CaptureException(err)
		return err
	}

	u.Salt = salt
	u.Password = HashPassword(salt, []byte(password))

	return nil
}

func HashPassword(salt, password []byte) []byte {
	hasher := sha512.New384()
	hasher.Write(salt)
	hasher.Write([]byte(password))

	return hasher.Sum(nil)
}
