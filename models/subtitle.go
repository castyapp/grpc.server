package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subtitle struct {
	ID            *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MediaSourceID *primitive.ObjectID `bson:"media_source_id,omitempty" json:"media_source_id,omitempty"`
	Lang          string              `bson:"lang" json:"size"`
	File          string              `bson:"file" json:"file"`
	CreatedAt     time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt     time.Time           `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
