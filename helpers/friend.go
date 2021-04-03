package helpers

import (
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/models"
	"github.com/golang/protobuf/ptypes"
)

func NewProtoFriend(friend *models.Friend) (*proto.Friend, error) {
	createdAt,  _ := ptypes.TimestampProto(friend.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(friend.UpdatedAt)
	return &proto.Friend{
		Accepted:  friend.Accepted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
