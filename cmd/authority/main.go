package main

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jessevdk/go-flags"

	"github.com/Quard/authority/internal/internal_api"
	"github.com/Quard/authority/internal/rest_api"
	"github.com/Quard/authority/internal/storage"
)

var opts struct {
	Bind            string `short:"b" long:"bind" env:"BIND" default:"localhost:5000" description:"address:port to listen"`
	InternalAPIBind string `long:"internal-api-bind" env:"INTERNAL_API_BIND" default:"localhost:5001"`
	MongoURI        string `long:"mongo-uri" env:"MONGO_URI" default:"mongodb://localhost:27017/authority"`
	SentryDSN       string `long:"sentry-dsn" env:"SENTRY_DSN"`
}

func main() {
	parser := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	sentry.Init(sentry.ClientOptions{
		Dsn: opts.SentryDSN,
	})
	defer sentry.Flush(time.Second * 5)

	stor := storage.NewMongoStorage(opts.MongoURI)

	restAPI := rest_api.NewRestAPIServer(opts.Bind, stor)
	internalAPI := internal_api.NewInternalAPIServer(opts.InternalAPIBind, stor)

	go restAPI.Run()
	go internalAPI.Run()

	select {}
}
