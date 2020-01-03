package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Friend struct {
	ID         *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`

	FriendId   *primitive.ObjectID   `bson:"friend_id, omitempty" json:"friend_id, omitempty"`
	UserId     *primitive.ObjectID   `bson:"user_id, omitempty" json:"user_id, omitempty"`

	Accepted   bool                  `bson:"accepted, omitempty" json:"accepted, omitempty"`

	CreatedAt  time.Time             `bson:"created_at, omitempty" json:"created_at, omitempty"`
	UpdatedAt  time.Time             `bson:"updated_at, omitempty" json:"updated_at, omitempty"`
}
