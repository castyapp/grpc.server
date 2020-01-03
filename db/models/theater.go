package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Theater struct {
	ID                 *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`
	Title              string                `bson:"title, omitempty" json:"title, omitempty"`

	// Theater hash. this will use for websocket connections
	Hash               string  `bson:"hash, omitempty" json:"hash, omitempty"`

	// 0 is just for the user
	// 1 is for everyone
	// 2 is for friends
	Privacy            int     `bson:"privacy, omitempty" json:"privacy, omitempty"`

	// 0 is just for the user
	// 1 is for everyone
	// 2 is for friends
	VideoPlayerAccess  int     `bson:"video_player_access, omitempty" json:"video_player_access, omitempty"`

	UserId             *primitive.ObjectID   `bson:"user_id, omitempty" json:"user_id, omitempty"`
	Movie              Movie                 `bson:"movie, omitempty" json:"movie, omitempty"`

	CreatedAt          time.Time    `bson:"created_at, omitempty" json:"created_at, omitempty"`
	UpdatedAt          time.Time    `bson:"updated_at, omitempty" json:"updated_at, omitempty"`
}

type TheaterMember struct {
	ID           *primitive.ObjectID  `bson:"_id, omitempty" json:"id, omitempty"`
	UserId       *primitive.ObjectID  `bson:"user_id, omitempty" json:"user_id, omitempty"`
	TheaterId    *primitive.ObjectID  `bson:"theater_id, omitempty" json:"theater_id, omitempty"`
	CreatedAt    time.Time            `bson:"created_at, omitempty" json:"created_at, omitempty"`
}