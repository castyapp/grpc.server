package db

import (
	"context"
	"fmt"
	"github.com/CastyLab/grpc.server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	Connection *mongo.Database
)

func Configure() error {

	ctx, _ := context.WithTimeout(context.Background(), 20 * time.Second)

	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.Map.Secrets.Db.Host, config.Map.Secrets.Db.Port))
	opts.SetAuth(options.Credential{
		Username: config.Map.Secrets.Db.User,
		Password: config.Map.Secrets.Db.Pass,
		AuthSource: config.Map.Secrets.Db.Name,
	})

	client, err := mongo.NewClient(opts)
	if err != nil {
		return fmt.Errorf("could not create new mongodb client: %v", err)
	}

	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("could not connect to mongodb client: %v", err)
	}

	Connection = client.Database(config.Map.Secrets.Db.Name)
	return nil
}