package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type Service struct {}

func (s *Service) RemoveActivity(ctx context.Context, req *proto.AuthenticateRequest) (*proto.Response, error) {

	user, err := auth.Authenticate(req)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	mdCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	var (
		filter = bson.M{"_id": user.ID}
		update = bson.M{
			"$set": bson.M{
				"activity": bson.M{},
			},
		}
	)

	if _, err := db.Connection.Collection("users").UpdateOne(mdCtx, filter, update); err != nil {
		sentry.CaptureException(err)
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "The requested parameter is not updated!",
		}, nil
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
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	mdCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	activityObjectId, err := primitive.ObjectIDFromHex(req.Activity.Id)
	if err != nil {

		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusNotAcceptable,
			Message: "Activity id is invalid!",
		}, nil
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

	if _, err := db.Connection.Collection("users").UpdateOne(mdCtx, filter, update); err != nil {
		sentry.CaptureException(err)
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "The requested parameter is not updated!",
		}, nil
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "The requested parameter is updated successfully!",
	}, nil
}

func (s *Service) UpdateState(ctx context.Context, req *proto.UpdateStateRequest) (*proto.Response, error) {

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	mdCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	var (
		filter = bson.M{"_id": user.ID}
		update = bson.M{
			"$set": bson.M{
				"state": int(req.State),
			},
		}
	)

	if _, err := db.Connection.Collection("users").UpdateOne(mdCtx, filter, update); err != nil {
		sentry.CaptureException(err)
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "The requested parameter is not updated!",
		}, nil
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "The requested parameter is updated successfully!",
	}, nil
}

func SetDBUserToProtoUser(user *models.User) (*proto.User, error) {

	lastLogin, _ := ptypes.TimestampProto(user.LastLogin)
	joinedAt,  _ := ptypes.TimestampProto(user.JoinedAt)
	updatedAt, _ := ptypes.TimestampProto(user.UpdatedAt)

	protoUser := &proto.User{
		Id:             user.ID.Hex(),
		Fullname:       user.Fullname,
		Username:       user.Username,
		Hash:           user.Hash,
		Email:          user.Email,
		IsActive:       user.IsActive,
		IsStaff:        user.IsStaff,
		Verified:       user.Verified,
		EmailVerified:  user.EmailVerified,
		Avatar:         user.Avatar,
		State:          proto.PERSONAL_STATE(user.State),
		LastLogin:      lastLogin,
		JoinedAt:       joinedAt,
		UpdatedAt:      updatedAt,
	}

	if user.Activity.ID != nil {
		protoUser.Activity = &proto.Activity{
			Id: user.Activity.ID.Hex(),
			Activity: user.Activity.Activity,
		}
	}

	return protoUser, nil
}

func (s *Service) GetUser(ctx context.Context, req *proto.AuthenticateRequest) (*proto.GetUserResponse, error) {

	user, err := auth.Authenticate(req)
	if err != nil {
		return nil, err
	}

	protoUser, err := SetDBUserToProtoUser(user)
	if err != nil {
		sentry.CaptureException(err)
		return &proto.GetUserResponse{
			Message: "Could not decode user!",
			Status: "failed",
			Code:   http.StatusInternalServerError,
		}, nil
	}

	return &proto.GetUserResponse{
		Result: protoUser,
		Status: "success",
		Code:   http.StatusOK,
	}, nil
}