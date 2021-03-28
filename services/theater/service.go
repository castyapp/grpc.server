package theater

import (
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/core"
)

type Service struct {
	*core.Context
	proto.UnimplementedTheaterServiceServer
}

func NewService(ctx *core.Context) *Service {
	return &Service{Context: ctx}
}
