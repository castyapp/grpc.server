package theater

import (
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/core"
)

type Service struct {
	*core.Context
	proto.UnimplementedTheaterServiceServer
}

func NewService(ctx *core.Context) *Service {
	return &Service{Context: ctx}
}
