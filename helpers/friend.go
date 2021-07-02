package helpers

import (
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewProtoFriend(f *models.Friend) (*proto.Friend, error) {
	return &proto.Friend{
		Accepted:  f.Accepted,
		CreatedAt: timestamppb.New(f.CreatedAt),
		UpdatedAt: timestamppb.New(f.UpdatedAt),
	}, nil
}
