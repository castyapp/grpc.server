package redis

import (
	"context"
	"fmt"
	"github.com/CastyLab/grpc.server/config"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	Client *redis.Client
)

func Configure() error {
	Client = redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: config.Map.Secrets.Redis.Sentinels,
		Password: config.Map.Secrets.Redis.Pass,
		MasterName: config.Map.Secrets.Redis.MasterName,
		DB: 0,
	})
	cmd := Client.Ping(context.Background())
	if res := cmd.Val(); res != "PONG" {
		log.Println("SentinelAddrs: ", config.Map.Secrets.Redis.Sentinels)
		mErr := fmt.Errorf("could not ping the redis server: %v", cmd.Err())
		sentry.CaptureException(mErr)
		return mErr
	}
	return nil
}

func Close() error {
	return Client.Close()
}
