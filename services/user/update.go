package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return &proto.GetUserResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "User updated successfully!",
			Result:  protoUser,
		}, nil
	}

	return nil, failedResponse
}