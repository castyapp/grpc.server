package helpers

import (
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/golang/protobuf/ptypes"
)

func NewProtoConnection(conn *models.Connection) *proto.Connection {

	createdAt, _ := ptypes.TimestampProto(conn.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(conn.UpdatedAt)

	return &proto.Connection{
		Id:             conn.ID.Hex(),
		ServiceUserId:  conn.ServiceUserId,
		Name:           conn.Name,
		Type:           conn.Type,
		AccessToken:    conn.AccessToken,
		ShowActivity:   conn.ShowActivity,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
