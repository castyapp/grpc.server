package providers

import (
	"fmt"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
)

type ConfigProvider struct{}

func (p *ConfigProvider) Register(ctx *core.Context) error {
	configFilePath := ctx.MustGetString("config.filepath")
	configMap, err := config.LoadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("could not load config: %v", err)
	}
	ctx.Set("config.map", configMap)
	return nil
}

func (p *ConfigProvider) Close(ctx *core.Context) error {
	return nil
}
