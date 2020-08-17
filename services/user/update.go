package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/internal"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

func (s *Service) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.GetUserResponse, error) {

	var (
		collection     = db.Connection.Collection("users")
		failedResponse = status.Error(codes.Internal, "Could not update the user, Please try again later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": user.ID}
	setUpdate := bson.M{}

	if req.Result.Fullname != "" && user.Fullname != req.Result.Fullname {
		setUpdate["fullname"] = req.Result.Fullname
	}

	if req.Result.Avatar != "" && user.Avatar != req.Result.Avatar {
		setUpdate["avatar"] = req.Result.Avatar
	}

	if len(setUpdate) == 0 {
		protoUser, err := helpers.NewProtoUser(user)
		if err != nil {
			sentry.CaptureException(err)
			return nil, status.Error(codes.Internal, "Internal server error!")
		}
		return &proto.GetUserResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "User updated successfully!",
			Result:  protoUser,
		}, nil
	}

	update := bson.M{"$set": setUpdate}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, failedResponse
	}

	dbUpdatedUser := new(models.User)
	if err := collection.FindOne(ctx, filter).Decode(dbUpdatedUser); err != nil {
		return nil, failedResponse
	}

	protoUser, err := helpers.NewProtoUser(dbUpdatedUser)
	if err != nil {
		return nil, failedResponse
	}

	if result.ModifiedCount != 0 {

		// sending updated user to websocket
		err := internal.Client.UserService.SendUpdateUserEvent(req.AuthRequest, protoUser.Id)
		if err != nil {
			sentry.CaptureException(err)
		}

		return &proto.GetUserResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "User updated successfully!",
			Result:  protoUser,
		}, nil
	}

	return nil, failedResponse
}

func (s *Service) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordRequest) (*proto.Response, error) {

	var (
		collection     = db.Connection.Collection("users")
		failedResponse = status.Error(codes.Internal, "Could not update the user's password, Please try again later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if !auth.ValidatePassword(user, req.CurrentPassword) {
		return nil, status.Error(codes.InvalidArgument, "Invalid Credentials!")
	}

	if req.NewPassword != req.VerifyNewPassword {
		return nil, status.Error(codes.InvalidArgument, "Passwords does not match!")
	}

	log.Println(models.HashPassword(req.VerifyNewPassword))
	log.Println(req.VerifyNewPassword)

	var (
		filter    = bson.M{"_id": user.ID}
		update    = bson.M{
			"$set": bson.M{
				"password": models.HashPassword(req.VerifyNewPassword),
			},
		}
	)

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, failedResponse
	}

	if result.ModifiedCount != 0 {
		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Password updated successfully!",
		}, nil
	}

	return nil, failedResponse
}