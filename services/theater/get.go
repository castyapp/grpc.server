package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *Service) GetTheater(ctx context.Context, theater *proto.Theater) (*proto.UserTheaterResponse, error) {

	var (
		collection     = db.Connection.Collection("theaters")
		failedResponse = status.Error(codes.Internal, "Could not get theater, Please try again later!")
	)

	if theater.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "TheaterId is required")
	}

	objectId, _ := primitive.ObjectIDFromHex(theater.Id)

	var dbTheater = new(models.Theater)
	if err := collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(dbTheater); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater!")
	}

	theater, err := helpers.NewTheaterProto(ctx, dbTheater)
	if err != nil {
		return nil, failedResponse
	}

	return &proto.UserTheaterResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  theater,
	}, nil
}
