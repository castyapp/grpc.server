package models

import "github.com/CastyLab/grpc.proto/proto"

type Movie struct {
	Type             proto.MovieType  `bson:"type" json:"type"`
	Uri              string           `bson:"uri" json:"uri"`

	Poster           string           `bson:"poster" json:"poster"`
	Subtitles        []Subtitle       `bson:"subtitles" json:"subtitles"`

	Size             int              `bson:"size" json:"size"`
	Length           int              `bson:"length" json:"length"`
	LastPlayedTime   int64            `bson:"last_played_time" json:"last_played_time"`
}
