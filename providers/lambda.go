package providers

import "github.com/castyapp/grpc.server/core"

type LambdaProvider struct {
	Registeration func(ctx *core.Context) error
	Closing       func(ctx *core.Context) error
}

func (p *LambdaProvider) Register(ctx *core.Context) error {
	if p.Registeration != nil {
		return p.Registeration(ctx)
	}
	return nil
}

func (p *LambdaProvider) Close(ctx *core.Context) error {
	if p.Closing != nil {
		return p.Closing(ctx)
	}
	return nil
}
