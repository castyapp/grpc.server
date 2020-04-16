package helpers

import (
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/golang/protobuf/ptypes"
)

func SetFriendToProto(friend *models.Friend) (*proto.Friend, error) {

	createdAt,  _ := ptypes.TimestampProto(friend.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(friend.UpdatedAt)

	protoUser := &proto.Friend{
		Accepted:  friend.Accepted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return protoUser, nil

}
