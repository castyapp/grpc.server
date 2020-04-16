package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func (s *Service) RollbackStates(ctx context.Context, req *proto.RollbackStatesRequest) (*proto.Response, error) {

	var (
		database = db.Connection
		collection = database.Collection("users")
		mCtx, _ = context.WithTimeout(ctx, 10 * time.Second)
	)

	update := bson.M{
		"$set": bson.M{
			"state": int(proto.PERSONAL_STATE_OFFLINE),
		},
	}

	filter := bson.M{
		"$or": []interface{} {
			bson.M{"state": int(proto.PERSONAL_STATE_ONLINE)},
			bson.M{"state": int(proto.PERSONAL_STATE_BUSY)},
			bson.M{"state": int(proto.PERSONAL_STATE_INVISIBLE)},
			bson.M{"state": int(proto.PERSONAL_STATE_IDLE)},
		},
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