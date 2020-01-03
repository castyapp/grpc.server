package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var (
	Client *mongo.Client
	Connection *mongo.Database
)

func init() {

	ctx, _ := context.WithTimeout(context.Background(), 20 * time.Second)

	var (
		host = os.Getenv("DB_HOST")
		port = os.Getenv("DB_PORT")
	)

	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))
	opts.SetAuth(options.Credential{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		PasswordSet: true,
	})

	client, err := mongo.NewClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	Client = client
	Connection = client.Database(os.Getenv("DB_NAME"))
}