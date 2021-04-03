package user

import (
	"context"
	"net/http"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
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

	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"fullname": "text",
			"username": "text",
		},
	})

	filter := bson.M{
		"_id": bson.M{"$ne": user.ID},
		"$text": bson.M{
			"$search": req.Keyword,
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
