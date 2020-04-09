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

func (s *Service) CreateMessage(ctx context.Context, req *proto.CreateMessageRequest) (*proto.CreateMessageResponse, error) {

	var (
		reciever         = new(models.User)
		collection       = db.Connection.Collection("messages")
		usersCollection  = db.Connection.Collection("users")
		failedResponse   = &proto.CreateMessageResponse{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not create message, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.CreateMessageResponse{
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

	if _, err := collection.InsertOne(ctx, message); err != nil {
		return failedResponse, err
	}

	protoUser, _ := helpers.SetDBUserToProtoUser(user)
	protoReciever, _ := helpers.SetDBUserToProtoUser(reciever)
	nowTime, _ := ptypes.TimestampProto(time.Now())

	return &proto.CreateMessageResponse{
		Code: http.StatusOK,
		Status: "success",
		Message: "Message created successfully!",
		Result:  &proto.Message{
			Content: req.Content,
			Sender:  protoUser,
			Reciever: protoReciever,
			Edited:  false,
			Deleted: false,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		},
	}, nil
}