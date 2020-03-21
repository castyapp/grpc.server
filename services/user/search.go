package user

import (
	"context"
	"errors"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.proto/messages"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
	"time"
)

func (s *Service) Search(ctx context.Context, req *proto.SearchUserRequest) (*proto.SearchUserResponse,error) {

	var (
		mCtx, _ = context.WithTimeout(ctx, 20 * time.Second)
		collection = db.Connection.Collection("users")
		emptyResponse = &proto.SearchUserResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Result:  make([]*messages.User, 0),
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
		"$or": []interface {}{
			bson.M{"fullname": bson.M{"$regex": req.Keyword}},
			bson.M{"fullname": bson.M{"$regex": strings.ToLower(req.Keyword)}},
			bson.M{"username": bson.M{"$regex": req.Keyword}},
			bson.M{"username": bson.M{"$regex": strings.ToLower(req.Keyword)}},
		},
	}

	cursor, err := collection.Find(mCtx, filter)
	if err != nil {
		return emptyResponse, nil
	}

	var protoUsers []*messages.User
	for cursor.Next(mCtx) {
		var dbUser = new(models.User)
		if err := cursor.Decode(dbUser); err != nil {
			break
		}
		protoUser, err := SetDBUserToProtoUser(dbUser)
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
