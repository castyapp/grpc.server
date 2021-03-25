package message

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.proto/protocol"
	"github.com/castyapp/grpc.server/db"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {

	var (
		reciever        = new(models.User)
		collection      = db.Connection.Collection("messages")
		usersCollection = db.Connection.Collection("users")
		failedResponse  = status.Error(codes.Internal, "Could not create message, Please try again later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if user.Username == req.Message.Reciever.Id {
		return nil, errors.New("receiver can not be you")
	}

	if err := usersCollection.FindOne(ctx, bson.M{"username": req.Message.Reciever.Id}).Decode(reciever); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find reciever!")
	}

	message := bson.M{
		"content":     req.Message.Content,
		"sender_id":   user.ID,
		"receiver_id": reciever.ID,
		"edited":      false,
		"deleted":     false,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	if _, err := collection.InsertOne(ctx, message); err != nil {
		return nil, failedResponse
	}

	nowTime, _ := ptypes.TimestampProto(time.Now())
	protoMessage := &proto.Message{
		Content:   req.Message.Content,
		Sender:    helpers.NewProtoUser(user),
		Reciever:  helpers.NewProtoUser(reciever),
		Edited:    false,
		Deleted:   false,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	buffer, err := protocol.NewMsgProtobuf(proto.EMSG_CHAT_MESSAGES, &proto.ChatMsgEvent{
		Message:   []byte(protoMessage.Content),
		Sender:    protoMessage.Sender,
		Reciever:  protoMessage.Reciever,
		CreatedAt: protoMessage.CreatedAt,
	})
	if err == nil {
		helpers.SendEventToUsers(ctx, buffer.Bytes(), []*proto.User{
			protoMessage.Sender,
			protoMessage.Reciever,
		})
	}

	return &proto.MessageResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Message created successfully!",
		Result:  protoMessage,
	}, nil
}
