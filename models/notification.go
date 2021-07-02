package models

import (
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID         *primitive.ObjectID                  `bson:"_id,omitempty" json:"id,omitempty"`
	Type       proto.Notification_NOTIFICATION_TYPE `bson:"type,omitempty" json:"type,omitempty"`
	Extra      *primitive.ObjectID                  `bson:"extra,omitempty" json:"extra,omitempty"`
	Read       bool                                 `bson:"read,omitempty" json:"read,omitempty"`
	FromUserID *primitive.ObjectID                  `bson:"from_user_id,omitempty" json:"from,omitempty"`
	ToUserID   *primitive.ObjectID                  `bson:"to_user_id,omitempty" json:"to,omitempty"`
	ReadAt     time.Time                            `bson:"read_at,omitempty" json:"read_at,omitempty"`
	CreatedAt  time.Time                            `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt  time.Time                            `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
