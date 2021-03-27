package message

import (
	"context"
	"errors"
	"net/http"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	*core.Context
	proto.UnimplementedMessagesServiceServer
}

func NewService(ctx *core.Context) *Service {
	return &Service{Context: ctx}
}

func (s *Service) GetUserMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {

	var (
		db              = s.MustGet("db.mongo").(*mongo.Database)
		reciever        = new(models.User)
		collection      = db.Collection("messages")
		usersCollection = db.Collection("users")
		failedResponse  = status.Error(codes.Internal, "Could not get messages, Please try again later!")
	)

	u, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if u.Username == req.ReceiverId {
		return nil, errors.New("receiver can not be you")
	}

	if err := usersCollection.FindOne(ctx, bson.M{"username": req.ReceiverId}).Decode(reciever); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find receiver!")
	}

	filter := bson.M{
		"$or": []interface{}{
			bson.M{
				"sender_id":   u.ID,
				"receiver_id": reciever.ID,
			},
			bson.M{
				"receiver_id": u.ID,
				"sender_id":   reciever.ID,
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
		protoMessage, err := helpers.NewProtoMessage(db, ctx, message)
		if err != nil {
			continue
		}
		protoMessages = append(protoMessages, protoMessage)
	}

	return &proto.GetMessagesResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: protoMessages,
	}, nil
}
