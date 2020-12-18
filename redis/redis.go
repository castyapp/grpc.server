package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/CastyLab/grpc.server/config"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
)

func Configure() error {

	if config.Map.Secrets.Redis.Replicaset {
		Client = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs:    config.Map.Secrets.Redis.Sentinels,
			SentinelPassword: config.Map.Secrets.Redis.SentinelPass,
			Password:         config.Map.Secrets.Redis.Pass,
			MasterName:       config.Map.Secrets.Redis.MasterName,
			DB:               0,
		})
	} else {
		Client = redis.NewClient(&redis.Options{
			Addr:     config.Map.Secrets.Redis.Addr,
			Password: config.Map.Secrets.Redis.Pass,
		})
	}

	cmd := Client.Ping(context.Background())
	if res := cmd.Val(); res != "PONG" {

		if config.Map.Secrets.Redis.Replicaset {
			log.Println("SentinelAddrs: ", config.Map.Secrets.Redis.Sentinels)
		} else {
			log.Println("Addr: ", config.Map.Secrets.Redis.Sentinels)
		}

		mErr := fmt.Errorf("could not ping the redis server: %v", cmd.Err())
		sentry.CaptureException(mErr)
		return mErr
	}
	return nil
}

func Close() error {
	return Client.Close()
}
