package user

import (
	"context"
	"errors"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (s *Service) RollbackStates(ctx context.Context, req *proto.RollbackStatesRequest) (*proto.Response, error) {

	var (
		database = db.Connection
		collection = database.Collection("users")
		mCtx, _ = context.WithTimeout(ctx, 10 * time.Second)
		update = bson.M{"$set": bson.M{"state": int(proto.PERSONAL_STATE_OFFLINE)}}
	)

	filter := bson.D{}
	for _, uId := range req.UsersIds {
		uObjectId, err := primitive.ObjectIDFromHex(uId)
		if err != nil {
			continue
		}
		filter = append(filter, bson.E{
			Key: "_id",
			Value: uObjectId,
		})
	}

	if len(filter) == 0 {
		return &proto.Response{
			Status: "Failed",
			Code: http.StatusInternalServerError,
		}, errors.New("users ids are empty")
	}

	if _, err := collection.UpdateMany(mCtx, filter, update); err != nil {
		return &proto.Response{
			Status: "Failed",
			Code: http.StatusInternalServerError,
		}, err
	}

	return &proto.Response{
		Status: "Success",
		Code: http.StatusOK,
	}, nil
}