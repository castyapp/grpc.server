package theater

import (
	"context"
	"log"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	c  *config.ConfigMap
	db *mongo.Database
	proto.UnimplementedTheaterServiceServer
}

func NewService(ctx context.Context) *Service {
	database := ctx.Value("db")
	if database == nil {
		log.Panicln("db value is required in context!")
	}
	configMap := ctx.Value("cm")
	if configMap == nil {
		log.Panicln("configMap value is required in context!")
	}
	return &Service{db: database.(*mongo.Database), c: configMap.(*config.ConfigMap)}
}
