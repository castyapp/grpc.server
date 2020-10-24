package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *Service) Search(ctx context.Context, req *proto.SearchUserRequest) (*proto.SearchUserResponse,error) {

	var (
		collection = db.Connection.Collection("users")
		emptyResponse = &proto.SearchUserResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Result:  make([]*proto.User, 0),
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if req.Keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "keyword is required")
	}

	filter := bson.M{
		"_id": bson.M{"$ne": user.ID},
		"$text": bson.M{"$search": req.Keyword},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return emptyResponse, nil
	}

	var protoUsers []*proto.User
	for cursor.Next(ctx) {
		var dbUser = new(models.User)
		if err := cursor.Decode(dbUser); err != nil {
			continue
		}
		protoUsers = append(protoUsers, helpers.NewProtoUser(dbUser))
	}

	return &proto.SearchUserResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  protoUsers,
	}, nil
}
