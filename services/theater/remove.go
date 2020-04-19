package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *Service) RemoveTheater(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.Response, error) {

	var (
		collection     = db.Connection.Collection("theaters")
		failedResponse = status.Error(codes.Internal, "Could not delete theater, Please try again later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if req.Theater.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "TheaterId is required")
	}

	theaterObjectId, err := primitive.ObjectIDFromHex(req.Theater.Id)
	if err != nil {
		return nil, failedResponse
	}

	filter := bson.M{
		"_id": theaterObjectId,
		"user_id": user.ID,
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, failedResponse
	}

	if result.DeletedCount == 1 {
		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Theater removed successfully!",
		}, nil
	}

	return nil, failedResponse
}