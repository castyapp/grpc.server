package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID            *primitive.ObjectID     `bson:"_id" json:"id"`

	Content       string                  `bson:"content" json:"content"`

	SenderId      *primitive.ObjectID     `bson:"sender_id" json:"sender_id"`
	ReceiverId    *primitive.ObjectID     `bson:"receiver_id" json:"receiver_id"`

	Edited        bool                    `bson:"edited" json:"edited"`
	Deleted       bool                    `bson:"deleted" json:"deleted"`

	CreatedAt     time.Time               `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time               `bson:"updated_at" json:"updated_at"`
	DeletedAt     time.Time               `bson:"deleted_at" json:"deleted_at"`
}