package message

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	c  *config.ConfigMap
	db *mongo.Database
	proto.UnimplementedMessagesServiceServer
}

func NewService(ctx context.Context) *Service {
	database := ctx.Value("db")
	if database == nil {
		log.Panicln("db value is required in context!")
	}
	configMap := ctx.Value("cm")
	if configMap == nil {
		log.Panicln("configMap value is required in context!")
	}
	return &Service{db: database.(*mongo.Database), c: configMap.(*config.ConfigMap)}
}

func (s *Service) GetUserMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {

	var (
		reciever        = new(models.User)
		collection      = s.db.Collection("messages")
		usersCollection = s.db.Collection("users")
		failedResponse  = status.Error(codes.Internal, "Could not get messages, Please try again later!")
	)

	u, err := auth.Authenticate(s.db, req.AuthRequest)
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
		protoMessage, err := helpers.NewProtoMessage(s.db, ctx, message)
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
