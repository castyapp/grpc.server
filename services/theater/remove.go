package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (s *Service) RemoveTheater(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.Response, error) {

	var (
		collection     = db.Connection.Collection("theaters")
		failedResponse = &proto.Response{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Could not remove theater, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	if req.Theater.Id == "" {
		return &proto.Response{
			Status:  "failed",
			Code:    420,
			Message: "Validation error, TheaterId is required!",
		}, nil
	}

	theaterObjectId, _ := primitive.ObjectIDFromHex(req.Theater.Id)

	filter := bson.M{
		"_id": theaterObjectId,
		"user_id": user.ID,
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return failedResponse, nil
	}

	if result.DeletedCount == 1 {
		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Theater removed successfully!",
		}, nil
	}

	return failedResponse, nil
}