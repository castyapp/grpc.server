package user

import (
	"context"
	"net/http"

	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Search(ctx context.Context, req *proto.SearchUserRequest) (*proto.SearchUserResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db            = dbConn.(*mongo.Database)
		collection    = db.Collection("users")
		emptyResponse = &proto.SearchUserResponse{
			Status: "success",
			Code:   http.StatusOK,
			Result: make([]*proto.User, 0),
		}
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if req.Keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "keyword is required")
	}

	filter := bson.M{
		"_id": bson.M{"$ne": user.ID},
		"$or": []interface{}{
			bson.M{
				"username": bson.M{
					"$regex":   req.Keyword,
					"$options": "i",
				},
			},
			bson.M{
				"fullname": bson.M{
					"$regex":   req.Keyword,
					"$options": "i",
				},
			},
		},
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
		Status: "success",
		Code:   http.StatusOK,
		Result: protoUsers,
	}, nil
}
