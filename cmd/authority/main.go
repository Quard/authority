package main

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jessevdk/go-flags"

	"github.com/Quard/authority/internal/rest_api"
	"github.com/Quard/authority/internal/storage"
)

var opts struct {
	Bind      string `short:"b" long:"bind" env:"BIND" default:"localhost:5001" description:"address:port to listen"`
	MongoURI  string `long:"mongo-uri" env:"MONGO_URI" default:"mongodb://localhost:27017/authority"`
	Secret    string `long:"secret" env:"SECRET"`
	SentryDSN string `long:"sentry-dsn" env:"SENTRY_DSN"`
}

func main() {
	parser := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	sentry.Init(sentry.ClientOptions{
		Dsn: opts.SentryDSN, // "https://fe3266fa658d4711bdeac482ebf4ddc4@sentry.io/1523280",
	})
	defer sentry.Flush(time.Second * 5)

	stor := storage.NewMongoStorage(opts.MongoURI)

	srv := rest_api.NewRestAPIServer(opts.Bind, stor)
	srv.Run()
}
