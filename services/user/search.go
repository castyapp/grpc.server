package user

import (
	"context"
	"errors"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func (s *Service) Search(ctx context.Context, req *proto.SearchUserRequest) (*proto.SearchUserResponse,error) {

	var (
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
		collection = db.Connection.Collection("users")
		emptyResponse = &proto.SearchUserResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Result:  make([]*proto.User, 0),
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.SearchUserResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	if req.Keyword == "" {
		return nil, errors.New("keyword is required")
	}

	filter := bson.M{
		"_id": bson.M{"$ne": user.ID},
		"$text": bson.M{"$search": req.Keyword},
	}

	cursor, err := collection.Find(mCtx, filter)
	if err != nil {
		return emptyResponse, nil
	}

	var protoUsers []*proto.User
	for cursor.Next(mCtx) {
		var dbUser = new(models.User)
		if err := cursor.Decode(dbUser); err != nil {
			break
		}
		protoUser, err := helpers.SetDBUserToProtoUser(dbUser)
		if err != nil {
			break
		}
		protoUsers = append(protoUsers, protoUser)
	}

	return &proto.SearchUserResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  protoUsers,
	}, nil
}
