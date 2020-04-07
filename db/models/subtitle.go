package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Subtitle struct {
	ID           *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`
	TheaterId    *primitive.ObjectID   `bson:"theater_id, omitempty" json:"theater_id, omitempty"`
	Lang         string                `bson:"lang" json:"size"`
	File         string                `bson:"file" json:"file"`
	CreatedAt    time.Time             `bson:"created_at, omitempty" json:"created_at, omitempty"`
	UpdatedAt    time.Time             `bson:"updated_at, omitempty" json:"updated_at, omitempty"`
}
