package message

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"movie.night.gRPC.server/db"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/proto"
	"movie.night.gRPC.server/services/auth"
	"net/http"
	"time"
)

func (s *Service) CreateMessage(ctx context.Context, req *proto.CreateMessageRequest) (*proto.Response, error) {

	var (
		reciever         = new(models.User)
		collection       = db.Connection.Collection("messages")
		usersCollection  = db.Connection.Collection("users")
		failedResponse   = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not create message, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, err
	}

	if user.Username == req.RecieverId {
		return failedResponse, errors.New("receiver can not be you")
	}

	if err := usersCollection.FindOne(ctx, bson.M{ "username": req.RecieverId }).Decode(reciever); err != nil {
		return failedResponse, err
	}

	message := bson.M{
		"content": req.Content,
		"sender_id": user.ID,
		"receiver_id": reciever.ID,
		"edited": false,
		"deleted": false,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	result, err := collection.InsertOne(ctx, message)
	if err != nil {
		return failedResponse, err
	}

	return &proto.Response{
		Code: http.StatusOK,
		Status: "success",
		Message: "Message created successfully!",
		Result:  []byte(result.InsertedID.(primitive.ObjectID).Hex()),
	}, nil
}