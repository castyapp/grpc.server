package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	ID         *primitive.ObjectID   `bson:"_id, omitempty" json:"id, omitempty"`
	Activity   string                `bson:"activity, omitempty" json:"activity, omitempty"`
}
