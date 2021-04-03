package models

import (
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaSource struct {
	ID        *primitive.ObjectID    `bson:"_id, omitempty" json:"id,omitempty"`
	UserId    *primitive.ObjectID    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Title     string                 `bson:"title" json:"title"`
	Type      proto.MediaSource_Type `bson:"type" json:"type,omitempty"`
	Banner    string                 `bson:"banner" json:"banner,omitempty"`
	Uri       string                 `bson:"uri" json:"uri,omitempty"`
	Length    int64                  `bson:"length" json:"length,omitempty"`
	Artist    string                 `bson:"artist" json:"artist,omitempty"`
	Subtitles []*Subtitle            `json:"subtitles,omitempty"`
	CreatedAt time.Time              `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
