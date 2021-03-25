package theater

import (
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/config"
)

type Service struct {
	c *config.ConfigMap
	proto.UnimplementedTheaterServiceServer
}
