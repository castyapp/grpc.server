package models

import (
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Theater struct {
	ID                *primitive.ObjectID       `bson:"_id,omitempty" json:"id,omitempty"`
	Description       string                    `bson:"description,omitempty" json:"description,omitempty"`
	Privacy           proto.PRIVACY             `bson:"privacy,omitempty" json:"privacy,omitempty"`
	VideoPlayerAccess proto.VIDEO_PLAYER_ACCESS `bson:"video_player_access,omitempty" json:"video_player_access,omitempty"`
	UserId            *primitive.ObjectID       `bson:"user_id,omitempty" json:"user_id,omitempty"`
	MediaSourceId     *primitive.ObjectID       `bson:"media_source_id,omitempty" json:"media_source_id,omitempty"`
	CreatedAt         time.Time                 `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt         time.Time                 `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Follow struct {
	ID                *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TheaterId         *primitive.ObjectID `bson:"theater_id,omitempty" json:"theater_id,omitempty"`
	UserId            *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	EmailNotification bool                `bson:"email_notification,omitempty" json:"email_notification,omitempty"`
	PushNotification  bool                `bson:"push_notification,omitempty" json:"push_notification,omitempty"`
	CreatedAt         time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt         time.Time           `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type TheaterMember struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TheaterId *primitive.ObjectID `bson:"theater_id,omitempty" json:"theater_id,omitempty"`
	UserId    *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
}
