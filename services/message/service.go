package message

import (
	"context"
	"errors"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type Service struct {}

func SetDbMessageToProtoMessage(ctx context.Context, message *models.Message) (*proto.Message, error) {

	var (
		dbSender   = new(models.User)
		collection = db.Connection.Collection("users")
	)

	if err := collection.FindOne(ctx, bson.M{ "_id": message.SenderId }).Decode(dbSender); err != nil {
		return nil, err
	}

	sender, err := helpers.SetDBUserToProtoUser(dbSender)
	if err != nil {
		return nil, err
	}

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

func (s *Service) GetUserMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {

	var (
		reciever         = new(models.User)
		collection       = db.Connection.Collection("messages")
		usersCollection  = db.Connection.Collection("users")
		failedResponse   = &proto.GetMessagesResponse{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not create message, Please try again later!",
		}
	)

	u, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.GetMessagesResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, err
	}

	if u.Username == req.ReceiverId {
		return failedResponse, errors.New("receiver can not be you")
	}

	mCtx, _ := context.WithTimeout(ctx, 10 * time.Second)

	if err := usersCollection.FindOne(mCtx, bson.M{ "username": req.ReceiverId }).Decode(reciever); err != nil {
		return failedResponse, err
	}

	filter := bson.M{
		"$or": []interface{} {
			bson.M{
				"sender_id": u.ID,
				"receiver_id": reciever.ID,
			},
			bson.M{
				"receiver_id": u.ID,
				"sender_id": reciever.ID,
			},
		},
	}

	cursor, err := collection.Find(mCtx, filter)
	if err != nil {
		return failedResponse, err
	}

	var protoMessages []*proto.Message
	for cursor.Next(mCtx) {
		var message = new(models.Message)
		if err := cursor.Decode(message); err != nil {
			break
		}
		protoMessage, err := SetDbMessageToProtoMessage(mCtx, message)
		if err != nil {
			break
		}
		protoMessages = append(protoMessages, protoMessage)
	}

	return &proto.GetMessagesResponse{
		Status: "success",
		Code: http.StatusOK,
		Result: protoMessages,
	}, nil
}