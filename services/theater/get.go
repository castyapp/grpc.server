package theater

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"movie.night.gRPC.server/db"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/proto"
	"movie.night.gRPC.server/services/auth"
	"net/http"
	"time"
)

func (s *Service) GetUserTheater(ctx context.Context, req *proto.GetTheaterRequest) (*proto.UserTheaterResponse, error) {

	var (
		collection     = db.Connection.Collection("theaters")
		failedResponse = &proto.UserTheaterResponse{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not get theater, Please try again later!",
		}
	)

	_, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.UserTheaterResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	if req.TheaterId == "" {
		return &proto.UserTheaterResponse{
			Status:  "failed",
			Code:    420,
			Message: "Validation error, TheaterId is required!",
		}, nil
	}

	objectId, _ := primitive.ObjectIDFromHex(req.TheaterId)

	filter := bson.M{
		"$or": []interface{} {
			bson.M{"hash": req.TheaterId},
			bson.M{"_id": objectId},
		},
	}

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	var dbTheater = new(models.Theater)
	if err := collection.FindOne(mCtx, filter).Decode(dbTheater); err != nil {
		return failedResponse, nil
	}

	theater, err := SetDbTheaterToMessageTheater(mCtx, dbTheater)
	if err != nil {
		return failedResponse, nil
	}

	return &proto.UserTheaterResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  theater,
	}, nil
}
