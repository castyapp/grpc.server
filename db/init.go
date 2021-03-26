package db

import (
	"context"
	"fmt"
	"time"

	"github.com/castyapp/grpc.server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Configure(c *config.ConfigMap) (*mongo.Database, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", c.DB.Host, c.DB.Port))
	opts.SetAuth(options.Credential{
		Username:   c.DB.User,
		Password:   c.DB.Pass,
		AuthSource: c.DB.Name,
	})

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("could not create new mongodb client: %v", err)
	}

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("could not connect to mongodb client: %v", err)
	}

	return client.Database(c.DB.Name), nil
}
