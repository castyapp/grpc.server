package helpers

import (
	"context"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewProtoMessage(db *mongo.Database, ctx context.Context, message *models.Message) (*proto.Message, error) {

	var (
		err        error
		sender     = new(models.User)
		collection = db.Collection("users")
	)

	if err := collection.FindOne(ctx, bson.M{"_id": message.SenderId}).Decode(sender); err != nil {
		return nil, err
	}

	createdAt, _ := ptypes.TimestampProto(message.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(message.UpdatedAt)

	protoMessage := &proto.Message{
		Id:        message.ID.Hex(),
		Content:   message.Content,
		Sender:    sender.ToProto(),
		Edited:    message.Edited,
		Deleted:   message.Deleted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if message.DeletedAt.Unix() != 0 {
		protoMessage.DeletedAt, err = ptypes.TimestampProto(message.DeletedAt)
		if err != nil {
			return nil, err
		}
	}

	return protoMessage, nil
}
