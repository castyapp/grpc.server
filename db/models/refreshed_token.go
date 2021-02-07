package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshedToken struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Token     string              `bson:"token,omitempty" json:"-,omitempty"`
	Valid     bool                `bson:"valid,omitempty" json:"valid,omitempty"`
	Csrf      string              `bson:"csrf,omitempty" json:"csrf,omitempty"`
	CreatedAt time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ExpiresAt time.Time           `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
}
