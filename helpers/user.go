package helpers

import (
	"context"
	"fmt"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/redis"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetFriendsFromDatabase(ctx context.Context, user *models.User) ([]*proto.User, error) {
	var (
		friends           = make([]*proto.User, 0)
		userCollection    = db.Connection.Collection("users")
		friendsCollection = db.Connection.Collection("friends")
	)

	filter := bson.M{
		"accepted": true,
		"$or": []interface{}{
			bson.M{"friend_id": user.ID},
			bson.M{"user_id": user.ID},
		},
	}

	cursor, err := friendsCollection.Find(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find friends!")
	}

	for cursor.Next(ctx) {

		var friend = new(models.Friend)
		if err := cursor.Decode(friend); err != nil {
			continue
		}

		var filter = bson.M{"_id": friend.FriendId}
		if user.ID.Hex() == friend.FriendId.Hex() {
			filter = bson.M{"_id": friend.UserId}
		}

		friendUserObject := new(models.UserWithState)
		if err := userCollection.FindOne(ctx, filter).Decode(friendUserObject); err != nil {
			continue
		}

		friends = append(friends, NewProtoUserWithState(friendUserObject))
	}

	return friends, nil
}

// update friends with new event of user
func SendEventToFriends(ctx context.Context, event []byte, user *models.User) error {
	friends, err := GetFriendsFromDatabase(ctx, user)
	if err != nil {
		return status.Error(codes.Internal, "Could not get friends!")
	}
	SendEventToUsers(ctx, event, friends)
	return nil
}

func SendEventToUser(ctx context.Context, event []byte, user *proto.User)  {
	redis.Client.Publish(ctx, fmt.Sprintf("user:events:%s", user.Id), event)
}

func SendEventToUsers(ctx context.Context, event []byte, users []*proto.User)  {
	for _, user := range users {
		redis.Client.Publish(ctx, fmt.Sprintf("user:events:%s", user.Id), event)
	}
}

func SendEventToTheaterMembers(ctx context.Context, event []byte, theater *models.Theater)  {
	redis.Client.Publish(ctx, fmt.Sprintf("theater:events:%s", theater.ID.Hex()), event)
}

func NewProtoUser(user *models.User) *proto.User {
	lastLogin, _ := ptypes.TimestampProto(user.LastLogin)
	joinedAt,  _ := ptypes.TimestampProto(user.JoinedAt)
	updatedAt, _ := ptypes.TimestampProto(user.UpdatedAt)
	return &proto.User{
		Id:             user.ID.Hex(),
		Fullname:       user.Fullname,
		Username:       user.Username,
		Hash:           user.Hash,
		Email:          user.Email,
		IsActive:       user.IsActive,
		IsStaff:        user.IsStaff,
		Verified:       user.Verified,
		EmailVerified:  user.EmailVerified,
		Avatar:         user.Avatar,
		TwoFaEnabled:   user.TwoFactorAuthEnabled,
		LastLogin:      lastLogin,
		JoinedAt:       joinedAt,
		UpdatedAt:      updatedAt,
	}
}

func NewProtoUserWithState(user *models.UserWithState) *proto.User {
	lastLogin, _ := ptypes.TimestampProto(user.LastLogin)
	joinedAt,  _ := ptypes.TimestampProto(user.JoinedAt)
	updatedAt, _ := ptypes.TimestampProto(user.UpdatedAt)
	return &proto.User{
		Id:             user.ID.Hex(),
		Fullname:       user.Fullname,
		Username:       user.Username,
		Hash:           user.Hash,
		Email:          user.Email,
		IsActive:       user.IsActive,
		IsStaff:        user.IsStaff,
		Verified:       user.Verified,
		EmailVerified:  user.EmailVerified,
		Avatar:         user.Avatar,
		TwoFaEnabled:   user.TwoFactorAuthEnabled,
		State:          user.State,
		LastLogin:      lastLogin,
		JoinedAt:       joinedAt,
		UpdatedAt:      updatedAt,
	}
}