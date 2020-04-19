package message

import (
	"context"
	"errors"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Service struct {}

func (s *Service) GetUserMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {

	var (
		reciever         = new(models.User)
		collection       = db.Connection.Collection("messages")
		usersCollection  = db.Connection.Collection("users")
		failedResponse   = status.Error(codes.Internal, "Could not get messages, Please try again later!")
	)

	u, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if u.Username == req.ReceiverId {
		return nil, errors.New("receiver can not be you")
	}

	if err := usersCollection.FindOne(ctx, bson.M{ "username": req.ReceiverId }).Decode(reciever); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find receiver!")
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

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, failedResponse
	}

	var protoMessages []*proto.Message
	for cursor.Next(ctx) {
		var message = new(models.Message)
		if err := cursor.Decode(message); err != nil {
			continue
		}
		protoMessage, err := helpers.NewProtoMessage(ctx, message)
		if err != nil {
			continue
		}
		protoMessages = append(protoMessages, protoMessage)
	}

	return &proto.GetMessagesResponse{
		Status: "success",
		Code: http.StatusOK,
		Result: protoMessages,
	}, nil
}