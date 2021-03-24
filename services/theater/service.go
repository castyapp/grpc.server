package theater

import (
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/config"
)

type Service struct {
	c *config.ConfMap
	proto.UnimplementedTheaterServiceServer
}
