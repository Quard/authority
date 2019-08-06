package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	uri  string
	conn *mongo.Client
}

func NewMongoStorage(uri string) MongoStorage {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return MongoStorage{
		uri:  uri,
		conn: client,
	}
}
