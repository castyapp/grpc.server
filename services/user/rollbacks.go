package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func (s *Service) RollbackStates(ctx context.Context, req *proto.RollbackStatesRequest) (*proto.Response, error) {

	var (
		database = db.Connection
		collection = database.Collection("users")
	)

	update := bson.M{
		"$set": bson.M{
			"state": int(proto.PERSONAL_STATE_OFFLINE),
		},
	}

	// filter non online users
	filter := bson.M{
		"state": bson.M{
			"$ne": int(proto.PERSONAL_STATE_OFFLINE),
		},
	}

	if _, err := collection.UpdateMany(ctx, filter, update); err != nil {
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