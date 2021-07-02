package helpers

import (
	"context"

	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewProtoMessage(ctx context.Context, db *mongo.Database, message *models.Message) (*proto.Message, error) {

	var (
		sender     = new(models.User)
		collection = db.Collection("users")
	)

	if err := collection.FindOne(ctx, bson.M{"_id": message.SenderID}).Decode(sender); err != nil {
		return nil, err
	}

	protoMessage := &proto.Message{
		Id:        message.ID.Hex(),
		Content:   message.Content,
		Sender:    sender.ToProto(),
		Edited:    message.Edited,
		Deleted:   message.Deleted,
		CreatedAt: timestamppb.New(message.CreatedAt),
		UpdatedAt: timestamppb.New(message.UpdatedAt),
	}

	if message.DeletedAt.Unix() != 0 {
		protoMessage.DeletedAt = timestamppb.New(message.DeletedAt)
	}

	return protoMessage, nil
}
