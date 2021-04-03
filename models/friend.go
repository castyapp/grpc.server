package models

import (
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Friend struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FriendId  *primitive.ObjectID `bson:"friend_id,omitempty" json:"friend_id,omitempty"`
	UserId    *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Accepted  bool                `bson:"accepted,omitempty" json:"accepted,omitempty"`
	CreatedAt time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time           `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (f *Friend) ToProto() *proto.Friend {
	createdAt, _ := ptypes.TimestampProto(f.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(f.UpdatedAt)
	return &proto.Friend{
		Accepted:  f.Accepted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
