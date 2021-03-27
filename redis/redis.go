package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
	"github.com/go-redis/redis/v8"
)

func Provider(ctx *core.Context) error {

	var (
		client *redis.Client
		cm     = ctx.MustGet("config.map").(*config.ConfigMap)
	)

	if cm.Redis.Cluster {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs:    cm.Redis.Sentinels,
			SentinelPassword: cm.Redis.SentinelPass,
			Password:         cm.Redis.Pass,
			MasterName:       cm.Redis.MasterName,
			DB:               0,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     cm.Redis.Addr,
			Password: cm.Redis.Pass,
		})
	}

	cmd := client.Ping(context.Background())
	if res := cmd.Val(); res != "PONG" {
		if cm.Redis.Cluster {
			log.Println(fmt.Sprintf("Redis: SentinelAddrs [%s]", cm.Redis.Sentinels))
		} else {
			log.Println(fmt.Sprintf("Redis: Addr [%s]", cm.Redis.Addr))
		}
		return fmt.Errorf("could not ping the redis server: %v", cmd.Err())
	}

	return ctx.Set("redis.conn", client)
}
