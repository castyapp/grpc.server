package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/castyapp/grpc.server/config"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
)

func Configure(c *config.ConfigMap) error {

	if c.Redis.Cluster {
		Client = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs:    c.Redis.Sentinels,
			SentinelPassword: c.Redis.SentinelPass,
			Password:         c.Redis.Pass,
			MasterName:       c.Redis.MasterName,
			DB:               0,
		})
	} else {
		Client = redis.NewClient(&redis.Options{
			Addr:     c.Redis.Addr,
			Password: c.Redis.Pass,
		})
	}

	cmd := Client.Ping(context.Background())
	if res := cmd.Val(); res != "PONG" {

		if c.Redis.Cluster {
			log.Println(fmt.Sprintf("Redis: SentinelAddrs [%s]", c.Redis.Sentinels))
		} else {
			log.Println(fmt.Sprintf("Redis: Addr [%s]", c.Redis.Addr))
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
