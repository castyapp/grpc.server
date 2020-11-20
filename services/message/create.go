package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.proto/protocol"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func (s *Service) CreateMessage(ctx context.Context, req *proto.CreateMessageRequest) (*proto.CreateMessageResponse, error) {

	var (
		reciever         = new(models.User)
		collection       = db.Connection.Collection("messages")
		usersCollection  = db.Connection.Collection("users")
		failedResponse   = status.Error(codes.Internal, "Could not create message, Please try again later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if user.Username == req.RecieverId {
		return nil, errors.New("receiver can not be you")
	}

	if err := usersCollection.FindOne(ctx, bson.M{ "username": req.RecieverId }).Decode(reciever); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find reciever!")
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
		return nil, failedResponse
	}

	nowTime, _ := ptypes.TimestampProto(time.Now())
	protoMessage := &proto.Message{
		Content:   req.Content,
		Sender:    helpers.NewProtoUser(user),
		Reciever:  helpers.NewProtoUser(reciever),
		Edited:    false,
		Deleted:   false,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	fromUser, err := json.Marshal(protoMessage.Sender)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("could not marshal Sender :%v", err).Error())
	}

	toUser, err := json.Marshal(protoMessage.Reciever)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("could not marshal Reciever :%v", err).Error())
	}

	buffer, err := protocol.NewMsgProtobuf(proto.EMSG_CHAT_MESSAGES, &proto.ChatMsgEvent{
		Message:   []byte(protoMessage.Content),
		From:      string(fromUser),
		To:        string(toUser),
		CreatedAt: protoMessage.CreatedAt,
	})
	if err == nil {
		helpers.SendEventToUsers(ctx, buffer.Bytes(), []*proto.User{
			protoMessage.Sender,
			protoMessage.Reciever,
		})
	}

	return &proto.CreateMessageResponse{
		Code: http.StatusOK,
		Status: "success",
		Message: "Message created successfully!",
		Result:  protoMessage,
	}, nil
}