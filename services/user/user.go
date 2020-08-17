package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Service struct {
	db *mongo.Database
}

func (s *Service) RemoveActivity(ctx context.Context, req *proto.AuthenticateRequest) (*proto.Response, error) {

	user, err := auth.Authenticate(req)
	if err != nil {
		return nil, err
	}

	var (
		filter = bson.M{"_id": user.ID}
		update = bson.M{
			"$set": bson.M{
				"activity": bson.M{},
			},
		}
	)

	if _, err := db.Connection.Collection("users").UpdateOne(ctx, filter, update); err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Aborted, "The requested parameter is not updated!")
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "The requested parameter is updated successfully!",
	}, nil
}

func (s *Service) UpdateActivity(ctx context.Context, req *proto.UpdateActivityRequest) (*proto.Response, error) {

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	activityObjectId, err := primitive.ObjectIDFromHex(req.Activity.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Activity id is invalid!")
	}

	var (
		filter = bson.M{"_id": user.ID}
		update = bson.M{
			"$set": bson.M{
				"activity": bson.M{
					"_id": activityObjectId,
					"activity": req.Activity.Activity,
				},
			},
		}
	)

	if _, err := db.Connection.Collection("users").UpdateOne(ctx, filter, update); err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Aborted, "The requested parameter is not updated!")
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "The requested parameter is updated successfully!",
	}, nil
}

func (s *Service) GetUser(_ context.Context, req *proto.AuthenticateRequest) (*proto.GetUserResponse, error) {

	user, err := auth.Authenticate(req)
	if err != nil {
		return nil, err
	}

	protoUser, err := helpers.NewProtoUser(user)
	if err != nil {
		sentry.CaptureException(err)
		return nil, status.Error(codes.Internal, "Could not decode user!")
	}

	return &proto.GetUserResponse{
		Result: protoUser,
		Status: "success",
		Code:   http.StatusOK,
	}, nil
}