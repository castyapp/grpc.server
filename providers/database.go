package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseProvider struct {
	client *mongo.Client
}

func (p *DatabaseProvider) Register(ctx *core.Context) error {
	mCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	cm := ctx.MustGet("config.map").(*config.ConfigMap)

	opts := options.Client()
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", cm.DB.Host, cm.DB.Port))
	opts.SetAuth(options.Credential{
		Username: cm.DB.User,
		Password: cm.DB.Pass,
		//AuthSource: cm.DB.Name,
	})

	client, err := mongo.NewClient(opts)
	if err != nil {
		return fmt.Errorf("could not create new mongodb client: %v", err)
	}
	p.client = client

	if err := client.Connect(mCtx); err != nil {
		return fmt.Errorf("could not connect to mongodb client: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping mongodb client: %v", err)
	}

	conn := client.Database(cm.DB.Name)
	return ctx.Set("db.mongo", conn)
}

func (p *DatabaseProvider) Close(ctx *core.Context) error {
	return p.client.Disconnect(ctx)
}
