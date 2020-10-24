package helpers

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
)

func NewProtoMessage(ctx context.Context, message *models.Message) (*proto.Message, error) {

	var (
		err error
		dbSender   = new(models.User)
		collection = db.Connection.Collection("users")
	)

	if err := collection.FindOne(ctx, bson.M{ "_id": message.SenderId }).Decode(dbSender); err != nil {
		return nil, err
	}

	sender := NewProtoUser(dbSender)
	createdAt, _ := ptypes.TimestampProto(message.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(message.UpdatedAt)

	protoMessage := &proto.Message{
		Id:       message.ID.Hex(),
		Content:  message.Content,
		Sender:   sender,
		Edited:   message.Edited,
		Deleted:  message.Deleted,
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