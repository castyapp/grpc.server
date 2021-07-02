package user

import (
	"context"
	"log"
	"net/http"

	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/libcasty-protocol-go/protocol"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	*core.Context
	proto.UnimplementedUserServiceServer
}

func NewService(ctx *core.Context) *Service {
	return &Service{Context: ctx}
}

func (s *Service) UpdateState(ctx context.Context, req *proto.UpdateStateRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	db := dbConn.(*mongo.Database)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}
	protoUser := helpers.NewProtoUser(user)

	var (
		filter = bson.M{"_id": user.ID}
		update = bson.M{"$set": bson.M{"state": req.State}}
	)

	if _, err := db.Collection("users").UpdateOne(ctx, filter, update); err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Aborted, "The requested parameter is not updated!")
	}

	// update self user with new state to other clients
	pms := &proto.PersonalStateMsgEvent{State: req.State, User: protoUser}
	buffer, err := protocol.NewMsgProtobuf(proto.EMSG_SELF_PERSONAL_STATE_CHANGED, pms)
	if err == nil {
		if err := helpers.SendEventToUser(s.Context, buffer.Bytes(), protoUser); err != nil {
			log.Println(err)
		}
	}

	// update friends with new state of user
	if buffer, err := protocol.NewMsgProtobuf(proto.EMSG_PERSONAL_STATE_CHANGED, pms); err == nil {
		if err := helpers.SendEventToFriends(s.Context, buffer.Bytes(), user); err != nil {
			return nil, err
		}
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "The requested parameter is updated successfully!",
	}, nil
}

func (s *Service) RemoveActivity(ctx context.Context, req *proto.AuthenticateRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	db := dbConn.(*mongo.Database)

	user, err := auth.Authenticate(s.Context, req)
	if err != nil {
		return nil, err
	}
	protoUser := helpers.NewProtoUser(user)

	var (
		filter = bson.M{"_id": user.ID}
		update = bson.M{
			"$set": bson.M{
				"activity": bson.M{},
			},
		}
	)

	if _, err := db.Collection("users").UpdateOne(ctx, filter, update); err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Aborted, "The requested parameter is not updated!")
	}

	// update self user with new activity to other clients
	pms := &proto.PersonalActivityMsgEvent{Activity: &proto.Activity{}, User: protoUser}
	buffer, err := protocol.NewMsgProtobuf(proto.EMSG_SELF_PERSONAL_ACTIVITY_CHANGED, pms)
	if err == nil {
		if err := helpers.SendEventToUser(s.Context, buffer.Bytes(), protoUser); err != nil {
			log.Println(err)
		}
	}

	// update friends with new activity of user
	if buffer, err := protocol.NewMsgProtobuf(proto.EMSG_PERSONAL_ACTIVITY_CHANGED, pms); err == nil {
		if err := helpers.SendEventToFriends(s.Context, buffer.Bytes(), user); err != nil {
			return nil, err
		}
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "The requested parameter is updated successfully!",
	}, nil
}

func (s *Service) UpdateActivity(ctx context.Context, req *proto.UpdateActivityRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	db := dbConn.(*mongo.Database)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}
	protoUser := helpers.NewProtoUser(user)

	activityObjectID, err := primitive.ObjectIDFromHex(req.Activity.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Activity id is invalid!")
	}

	var (
		filter = bson.M{"_id": user.ID}
		update = bson.M{
			"$set": bson.M{
				"activity": bson.M{
					"_id":      activityObjectID,
					"activity": req.Activity.Activity,
				},
			},
		}
	)

	if _, err := db.Collection("users").UpdateOne(ctx, filter, update); err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Aborted, "The requested parameter is not updated!")
	}

	activity := &proto.Activity{
		Id:       activityObjectID.Hex(),
		Activity: req.Activity.Activity,
	}

	// update self user with new activity to other clients
	pms := &proto.PersonalActivityMsgEvent{Activity: activity, User: protoUser}
	buffer, err := protocol.NewMsgProtobuf(proto.EMSG_SELF_PERSONAL_ACTIVITY_CHANGED, pms)
	if err == nil {
		if err := helpers.SendEventToUser(s.Context, buffer.Bytes(), protoUser); err != nil {
			log.Println(err)
		}
	}

	// update friends with new activity of user
	if buffer, err := protocol.NewMsgProtobuf(proto.EMSG_PERSONAL_ACTIVITY_CHANGED, pms); err == nil {
		if err := helpers.SendEventToFriends(s.Context, buffer.Bytes(), user); err != nil {
			return nil, err
		}
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "The requested parameter is updated successfully!",
	}, nil
}

func (s *Service) GetUser(_ context.Context, req *proto.AuthenticateRequest) (*proto.GetUserResponse, error) {
	user, err := auth.Authenticate(s.Context, req)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserResponse{
		Result: helpers.NewProtoUser(user),
		Status: "success",
		Code:   http.StatusOK,
	}, nil
}
