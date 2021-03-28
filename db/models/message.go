package models

import (
	"context"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	ID         *primitive.ObjectID `bson:"_id" json:"id"`
	Content    string              `bson:"content" json:"content"`
	SenderId   *primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	ReceiverId *primitive.ObjectID `bson:"receiver_id" json:"receiver_id"`
	Edited     bool                `bson:"edited" json:"edited"`
	Deleted    bool                `bson:"deleted" json:"deleted"`
	CreatedAt  time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time           `bson:"updated_at" json:"updated_at"`
	DeletedAt  time.Time           `bson:"deleted_at" json:"deleted_at"`
}

func (m *Message) ToProto(db *mongo.Database) (*proto.Message, error) {

	var (
		ctx        = context.Background()
		err        error
		sender     = new(User)
		collection = db.Collection("users")
	)

	if err := collection.FindOne(ctx, bson.M{"_id": m.SenderId}).Decode(sender); err != nil {
		return nil, err
	}

	createdAt, _ := ptypes.TimestampProto(m.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(m.UpdatedAt)

	protoMessage := &proto.Message{
		Id:        m.ID.Hex(),
		Content:   m.Content,
		Sender:    sender.ToProto(),
		Edited:    m.Edited,
		Deleted:   m.Deleted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if m.DeletedAt.Unix() != 0 {
		protoMessage.DeletedAt, err = ptypes.TimestampProto(m.DeletedAt)
		if err != nil {
			return nil, err
		}
	}

	return protoMessage, nil
}
