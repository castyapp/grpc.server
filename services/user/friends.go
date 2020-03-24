package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func (s *Service) GetFriends(ctx context.Context, req *proto.AuthenticateRequest) (*proto.FriendsResponse, error) {

	var (
		friends []*proto.User

		database   = db.Connection
		mCtx, _    = context.WithTimeout(ctx, 20 * time.Second)

		userCollection    = database.Collection("users")
		friendsCollection = database.Collection("friends")

		failedResponse = &proto.FriendsResponse{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not get friends, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req)
	if err != nil {
		return &proto.FriendsResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	filter := bson.M{
		"accepted": true,
		"$or": []interface{}{
			bson.M{"friend_id": user.ID},
			bson.M{"user_id": user.ID},
		},
	}

	cursor, err := friendsCollection.Find(mCtx, filter)
	if err != nil {
		return failedResponse, nil
	}

	for cursor.Next(mCtx) {

		var friend = new(models.Friend)
		if err := cursor.Decode(friend); err != nil {
			break
		}

		var filter = bson.M{"_id": friend.FriendId}
		if user.ID.Hex() == friend.FriendId.Hex() {
			filter = bson.M{"_id": friend.UserId}
		}

		friendUserObject := new(models.User)
		if err := userCollection.FindOne(mCtx, filter).Decode(friendUserObject); err != nil {
			break
		}

		messageUser, err := SetDBUserToProtoUser(friendUserObject)
		if err != nil {
			break
		}

		friends = append(friends, messageUser)
	}

	return &proto.FriendsResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  friends,
	}, nil
}
