package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Quard/authority/internal/user"
)

type MongoStorage struct {
	uri  string
	conn *mongo.Database
}

func NewMongoStorage(uri string) MongoStorage {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return MongoStorage{
		uri:  uri,
		conn: client.Database("authority"),
	}
}

func (s MongoStorage) AddUser(user user.User) error {
	userCollection := s.conn.Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	document, errConv := convToInsertDocument(user)
	if errConv != nil {
		sentry.CaptureException(errConv)
		return errConv
	}
	_, err := userCollection.InsertOne(ctx, document)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	return nil
}

func (s MongoStorage) GetUserByEmail(email string) (user.User, error) {
	var user user.User
	userCollection := s.conn.Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	findResult := userCollection.FindOne(ctx, bson.M{"email": email})
	findResultErr := findResult.Err()
	if findResultErr != nil {
		sentry.CaptureException(findResultErr)
		return user, errors.New("couldn't retrieve user")
	}
	err := findResult.Decode(&user)
	if err != nil {
		return user, ErrUserNotFound
	}

	return user, nil
}

func convToInsertDocument(val interface{}) (bson.M, error) {
	var document bson.M
	bytes, err := bson.Marshal(val)
	if err != nil {
		sentry.CaptureException(err)
		return document, err
	}
	if err := bson.Unmarshal(bytes, &document); err != nil {
		sentry.CaptureException(err)
		return document, err
	}

	delete(document, "_id")

	return document, nil
}
