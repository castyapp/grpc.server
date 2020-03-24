package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func (s *Service) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.GetUserResponse, error) {

	var (
		database = db.Connection
		collection = database.Collection("users")
		failedResponse = &proto.GetUserResponse{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not update the user, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.GetUserResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	filter := bson.M{"_id": user.ID}
	setUpdate := bson.M{}

	if req.Result.Avatar != "" && user.Avatar != req.Result.Avatar {
		setUpdate["avatar"] = req.Result.Avatar
	}

	if req.Result.Fullname != "" && user.Fullname != req.Result.Fullname {
		setUpdate["fullname"] = req.Result.Fullname
	}

	if len(setUpdate) == 0 {
		return &proto.GetUserResponse{
			Status:  "failed",
			Code:    420,
			Message: "Fields are required!",
		}, nil
	}

	update := bson.M{"$set": setUpdate}
	result, err := collection.UpdateOne(mCtx, filter, update)
	if err != nil {
		return failedResponse, nil
	}

	dbUpdatedUser := new(models.User)
	if err := collection.FindOne(mCtx, filter).Decode(dbUpdatedUser); err != nil {
		return failedResponse, nil
	}

	protoUser, err := SetDBUserToProtoUser(dbUpdatedUser)
	if err != nil {
		return failedResponse, nil
	}

	if result.ModifiedCount != 0 {
		return &proto.GetUserResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "User updated successfully!",
			Result:  protoUser,
		}, nil
	}

	return failedResponse, nil
}