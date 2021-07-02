package providers

import (
	"fmt"
	"time"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
	"github.com/getsentry/sentry-go"
)

type SentryProvider struct{}

func (p *SentryProvider) Register(ctx *core.Context) error {
	cm := ctx.MustGet("config.map").(*config.Map)
	if cm.Sentry.Enabled {
		if err := sentry.Init(sentry.ClientOptions{Dsn: cm.Sentry.Dsn}); err != nil {
			return fmt.Errorf("could not initilize sentry: %v", err)
		}
	}
	return nil
}

func (p *SentryProvider) Close(ctx *core.Context) error {
	sentry.Flush(5 * time.Second)
	return nil
}
