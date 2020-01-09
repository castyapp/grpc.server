package user

import (
	"context"
	"gitlab.com/movienight1/grpc.proto"
	"go.mongodb.org/mongo-driver/bson"
	"movie.night.gRPC.server/db"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/services/auth"
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

	if err := userCollection.FindOne(mCtx, bson.M{ "username": string(req.FriendId) }).Decode(dbFriendUserObject); err != nil {
		return failedResponse, nil
	}

	filter := bson.M{
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

	friendUser, err := SetDBUserToProtoUser(dbFriendUserObject)
	if err != nil {
		return failedResponse, nil
	}

	return &proto.FriendResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  friendUser,
	}, nil
}
