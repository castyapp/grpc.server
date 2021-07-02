package models

import (
	"context"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	ID         *primitive.ObjectID `bson:"_id" json:"id"`
	Content    string              `bson:"content" json:"content"`
	SenderID   *primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	ReceiverID *primitive.ObjectID `bson:"receiver_id" json:"receiver_id"`
	Edited     bool                `bson:"edited" json:"edited"`
	Deleted    bool                `bson:"deleted" json:"deleted"`
	CreatedAt  time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time           `bson:"updated_at" json:"updated_at"`
	DeletedAt  time.Time           `bson:"deleted_at" json:"deleted_at"`
}

func (m *Message) ToProto(db *mongo.Database) (*proto.Message, error) {

	var (
		ctx        = context.Background()
		sender     = new(User)
		collection = db.Collection("users")
	)

	if err := collection.FindOne(ctx, bson.M{"_id": m.SenderID}).Decode(sender); err != nil {
		return nil, err
	}

	createdAt := timestamppb.New(m.CreatedAt)
	updatedAt := timestamppb.New(m.UpdatedAt)

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
		protoMessage.DeletedAt = timestamppb.New(m.DeletedAt)
	}

	return protoMessage, nil
}
