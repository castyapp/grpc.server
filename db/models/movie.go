package models

type Movie struct {
	MovieUri           string        `bson:"movie_uri" json:"movie_uri"`

	Poster             string        `bson:"poster" json:"poster"`
	Subtitles          []Subtitle    `bson:"subtitles" json:"subtitles"`

	Size               int           `bson:"size" json:"size"`
	Length             int           `bson:"length" json:"length"`
	LastPlayedTime     int64     `bson:"last_played_time" json:"last_played_time"`
}
