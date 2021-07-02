package helpers

import (
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewProtoConnection(c *models.Connection) *proto.Connection {
	return &proto.Connection{
		Id:            c.ID.Hex(),
		ServiceUserId: c.ServiceUserID,
		Name:          c.Name,
		Type:          c.Type,
		AccessToken:   c.AccessToken,
		ShowActivity:  c.ShowActivity,
		CreatedAt:     timestamppb.New(c.CreatedAt),
		UpdatedAt:     timestamppb.New(c.UpdatedAt),
	}
}
