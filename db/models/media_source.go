package models

import (
	"github.com/CastyLab/grpc.proto/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MediaSource struct {
	ID              *primitive.ObjectID       `bson:"_id, omitempty" json:"id, omitempty"`

	UserId          *primitive.ObjectID       `bson:"user_id, omitempty" json:"user_id, omitempty"`

	Title           string                    `bson:"title" json:"title"`
	Type            proto.MediaSource_Type    `bson:"type" json:"type,omitempty"`
	Banner          string                    `bson:"banner" json:"banner,omitempty"`
	Uri             string                    `bson:"uri" json:"uri,omitempty"`

	LastPlayedTime  int64                     `bson:"last_played_time" json:"last_played_time,omitempty"`
	Length          int64                     `bson:"length" json:"length,omitempty"`

	Subtitles       []*Subtitle               `json:"subtitles,omitempty"`

	CreatedAt       time.Time                 `bson:"created_at, omitempty" json:"created_at, omitempty"`
	UpdatedAt       time.Time                 `bson:"updated_at, omitempty" json:"updated_at, omitempty"`
}
