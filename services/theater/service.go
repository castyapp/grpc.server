package theater

import (
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/libcasty-protocol-go/proto"
)

type Service struct {
	*core.Context
	proto.UnimplementedTheaterServiceServer
}

func NewService(ctx *core.Context) *Service {
	return &Service{Context: ctx}
}
