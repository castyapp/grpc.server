package message

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/libcasty-protocol-go/protocol"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {

	var (
		db              = s.MustGet("db.mongo").(*mongo.Database)
		receiver        = new(models.User)
		collection      = db.Collection("messages")
		usersCollection = db.Collection("users")
		failedResponse  = status.Error(codes.Internal, "Could not create message, Please try again later!")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if user.Username == req.Message.Receiver.Id {
		return nil, errors.New("receiver can not be you")
	}

	if err := usersCollection.FindOne(ctx, bson.M{"username": req.Message.Receiver.Id}).Decode(receiver); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find receiver!")
	}

	message := bson.M{
		"content":     req.Message.Content,
		"sender_id":   user.ID,
		"receiver_id": receiver.ID,
		"edited":      false,
		"deleted":     false,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	if _, err := collection.InsertOne(ctx, message); err != nil {
		return nil, failedResponse
	}

	nowTime := timestamppb.New(time.Now())
	protoMessage := &proto.Message{
		Content:   req.Message.Content,
		Sender:    helpers.NewProtoUser(user),
		Receiver:  helpers.NewProtoUser(receiver),
		Edited:    false,
		Deleted:   false,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	buffer, err := protocol.NewMsgProtobuf(proto.EMSG_CHAT_MESSAGES, &proto.ChatMsgEvent{
		Message:   []byte(protoMessage.Content),
		Sender:    protoMessage.Sender,
		Receiver:  protoMessage.Receiver,
		CreatedAt: protoMessage.CreatedAt,
	})
	if err == nil {
		helpers.SendEventToUsers(s.Context, buffer.Bytes(), []*proto.User{
			protoMessage.Sender,
			protoMessage.Receiver,
		})
	}

	return &proto.MessageResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Message created successfully!",
		Result:  protoMessage,
	}, nil
}
