package user

import (
	"context"
	"net/http"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetFriend(ctx context.Context, req *proto.FriendRequest) (*proto.FriendResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                 = dbConn.(*mongo.Database)
		dbFriend           = new(models.Friend)
		dbFriendUserObject = new(models.User)
		userCollection     = db.Collection("users")
		friendsCollection  = db.Collection("friends")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	findFriendFilter := bson.M{"username": req.FriendId}

	friendObjectId, err := primitive.ObjectIDFromHex(req.FriendId)
	if err == nil {
		findFriendFilter = bson.M{"_id": friendObjectId}
	}

	if err := userCollection.FindOne(ctx, findFriendFilter).Decode(dbFriendUserObject); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find friend!")
	}

	filter := bson.M{
		"accepted": true,
		"$or": []interface{}{
			bson.M{
				"friend_id": user.ID,
				"user_id":   dbFriendUserObject.ID,
			},
			bson.M{
				"user_id":   user.ID,
				"friend_id": dbFriendUserObject.ID,
			},
		},
	}

	if err := friendsCollection.FindOne(ctx, filter).Decode(dbFriend); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find friend!")
	}

	return &proto.FriendResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: helpers.NewProtoUser(dbFriendUserObject),
	}, nil
}

func (s *Service) GetFriendRequest(ctx context.Context, req *proto.FriendRequest) (*proto.Friend, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                = dbConn.(*mongo.Database)
		dbFriend          = new(models.Friend)
		friendsCollection = db.Collection("friends")
		failedResponse    = &proto.Friend{}
		failedErr         = status.Error(codes.Internal, "Could not et friend request!")
	)

	if _, err := auth.Authenticate(s.Context, req.AuthRequest); err != nil {
		return nil, err
	}

	requestObjectId, err := primitive.ObjectIDFromHex(req.RequestId)
	if err != nil {
		return failedResponse, failedErr
	}

	if err := friendsCollection.FindOne(ctx, bson.M{"_id": requestObjectId}).Decode(dbFriend); err != nil {
		return failedResponse, status.Error(codes.NotFound, "Could not find friend request!")
	}

	friendUser, err := helpers.NewProtoFriend(dbFriend)
	if err != nil {
		return failedResponse, failedErr
	}

	return friendUser, nil
}
