package models

import (
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Connection struct {
	ID             *primitive.ObjectID   `bson:"_id" json:"id,omitempty"`
	ServiceUserId  string                `bson:"service_user_id" json:"service_user_id,omitempty"`
	Name           string                `bson:"name" json:"name,omitempty"`
	Type           proto.Connection_Type `bson:"type" json:"type,omitempty"`
	AccessToken    string                `bson:"access_token" json:"access_token,omitempty"`
	RefreshedToken string                `bson:"refreshed_token" json:"-"`
	ShowActivity   bool                  `bson:"show_activity,omitempty" json:"show_activity,omitempty"`
	UserId         *primitive.ObjectID   `bson:"user_id,omitempty" json:"user_id,omitempty"`
	CreatedAt      time.Time             `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      time.Time             `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Convert connection model to protobuf
func (c *Connection) ToProto() *proto.Connection {
	createdAt, _ := ptypes.TimestampProto(c.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(c.UpdatedAt)
	return &proto.Connection{
		Id:            c.ID.Hex(),
		ServiceUserId: c.ServiceUserId,
		Name:          c.Name,
		Type:          c.Type,
		AccessToken:   c.AccessToken,
		ShowActivity:  c.ShowActivity,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
