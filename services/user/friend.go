package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *Service) GetFriend(ctx context.Context, req *proto.FriendRequest) (*proto.FriendResponse, error) {

	var (
		dbFriend           = new(models.Friend)
		dbFriendUserObject = new(models.User)
		userCollection     = db.Connection.Collection("users")
		friendsCollection  = db.Connection.Collection("friends")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	findFriendFilter := bson.M{ "username": req.FriendId }

	friendObjectId, err := primitive.ObjectIDFromHex(req.FriendId)
	if err == nil {
		findFriendFilter = bson.M{ "_id": friendObjectId }
	}

	if err := userCollection.FindOne(ctx, findFriendFilter).Decode(dbFriendUserObject); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find friend!")
	}

	filter := bson.M{
		"accepted": true,
		"$or": []interface{}{
			bson.M{
				"friend_id": user.ID,
				"user_id": dbFriendUserObject.ID,
			},
			bson.M{
				"user_id": user.ID,
				"friend_id": dbFriendUserObject.ID,
			},
		},
	}

	if err := friendsCollection.FindOne(ctx, filter).Decode(dbFriend); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find friend!")
	}

	return &proto.FriendResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  helpers.NewProtoUser(dbFriendUserObject),
	}, nil
}

func (s *Service) GetFriendRequest(ctx context.Context, req *proto.FriendRequest) (*proto.Friend, error) {

	var (
		database   = db.Connection
		dbFriend   = new(models.Friend)
		friendsCollection = database.Collection("friends")
		failedResponse = &proto.Friend{}
		failedErr      = status.Error(codes.Internal, "Could not et friend request!")
	)

	if _, err := auth.Authenticate(req.AuthRequest); err != nil {
		return nil, err
	}

	requestObjectId, err := primitive.ObjectIDFromHex(req.RequestId)
	if err != nil {
		return failedResponse, failedErr
	}

	if err := friendsCollection.FindOne(ctx, bson.M{ "_id": requestObjectId }).Decode(dbFriend); err != nil {
		return failedResponse, status.Error(codes.NotFound, "Could not find friend request!")
	}

	friendUser, err := helpers.NewProtoFriend(dbFriend)
	if err != nil {
		return failedResponse, failedErr
	}

	return friendUser, nil
}
