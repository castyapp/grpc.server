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
	"net/http"
	"time"
)

func (s *Service) GetFriend(ctx context.Context, req *proto.FriendRequest) (*proto.FriendResponse, error) {

	var (

		database   = db.Connection

		dbFriend           = new(models.Friend)
		dbFriendUserObject = new(models.User)

		mCtx, _  = context.WithTimeout(ctx, 20 * time.Second)

		userCollection    = database.Collection("users")
		friendsCollection = database.Collection("friends")

		failedResponse = &proto.FriendResponse{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not get the friend, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.FriendResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	findFriendFilter := bson.M{ "username": req.FriendId }

	friendObjectId, err := primitive.ObjectIDFromHex(req.FriendId)
	if err == nil {
		findFriendFilter = bson.M{ "_id": friendObjectId }
	}

	if err := userCollection.FindOne(mCtx, findFriendFilter).Decode(dbFriendUserObject); err != nil {
		return failedResponse, nil
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

	if err := friendsCollection.FindOne(mCtx, filter).Decode(dbFriend); err != nil {
		return failedResponse, nil
	}

	friendUser, err := helpers.SetDBUserToProtoUser(dbFriendUserObject)
	if err != nil {
		return failedResponse, nil
	}

	return &proto.FriendResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  friendUser,
	}, nil
}

func (s *Service) GetFriendRequest(ctx context.Context, req *proto.FriendRequest) (*proto.Friend, error) {

	var (
		database   = db.Connection
		dbFriend   = new(models.Friend)
		mCtx, _    = context.WithTimeout(ctx, 20 * time.Second)
		friendsCollection = database.Collection("friends")
		failedResponse = &proto.Friend{}
	)

	if _, err := auth.Authenticate(req.AuthRequest); err != nil {
		return failedResponse, err
	}

	requestObjectId, err := primitive.ObjectIDFromHex(req.RequestId)
	if err != nil {
		return failedResponse, err
	}

	if err := friendsCollection.FindOne(mCtx, bson.M{ "_id": requestObjectId }).Decode(dbFriend); err != nil {
		return failedResponse, err
	}

	friendUser, err := helpers.SetFriendToProto(dbFriend)
	if err != nil {
		return failedResponse, err
	}

	return friendUser, nil
}
