package providers

import (
	"context"
	"fmt"
	"log"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
	"github.com/go-redis/redis/v8"
)

type RedisProvider struct {
	client *redis.Client
}

func (p *RedisProvider) Register(ctx *core.Context) error {
	cm := ctx.MustGet("config.map").(*config.ConfigMap)
	if cm.Redis.Cluster {
		p.client = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs:    cm.Redis.Sentinels,
			SentinelPassword: cm.Redis.SentinelPass,
			Password:         cm.Redis.Pass,
			MasterName:       cm.Redis.MasterName,
			DB:               0,
		})
	} else {
		p.client = redis.NewClient(&redis.Options{
			Addr:     cm.Redis.Addr,
			Password: cm.Redis.Pass,
		})
	}

	cmd := p.client.Ping(context.Background())
	if res := cmd.Val(); res != "PONG" {
		if cm.Redis.Cluster {
			log.Println(fmt.Sprintf("Redis: SentinelAddrs [%s]", cm.Redis.Sentinels))
		} else {
			log.Println(fmt.Sprintf("Redis: Addr [%s]", cm.Redis.Addr))
		}
		return fmt.Errorf("could not ping the redis server: %v", cmd.Err())
	}

	return ctx.Set("redis.conn", p.client)
}

func (p *RedisProvider) Close(ctx *core.Context) error {
	return p.client.Close()
}
