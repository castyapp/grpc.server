package helpers

import (
	"fmt"
	"log"

	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetFriendsFromDatabase(ctx *core.Context, user *models.User) ([]*proto.User, error) {

	db := ctx.MustGet("db.mongo").(*mongo.Database)

	var (
		friends           = make([]*proto.User, 0)
		userCollection    = db.Collection("users")
		friendsCollection = db.Collection("friends")
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

		var filter = bson.M{"_id": friend.FriendID}
		if user.ID.Hex() == friend.FriendID.Hex() {
			filter = bson.M{"_id": friend.UserID}
		}

		friendUserObject := new(models.User)
		if err := userCollection.FindOne(ctx, filter).Decode(friendUserObject); err != nil {
			continue
		}

		friends = append(friends, NewProtoUserWithState(friendUserObject))
	}

	return friends, nil
}

// update friends with new event of user
func SendEventToFriends(ctx *core.Context, event []byte, user *models.User) error {
	friends, err := GetFriendsFromDatabase(ctx, user)
	if err != nil {
		return status.Error(codes.Internal, "Could not get friends!")
	}
	SendEventToUsers(ctx, event, friends)
	return nil
}

func SendEventToUser(ctx *core.Context, event []byte, user *proto.User) (err error) {
	redisConn, err := ctx.Get("redis.conn")
	if err != nil {
		return err
	}
	_, err = redisConn.(*redis.Client).Publish(ctx, fmt.Sprintf("user:events:%s", user.Id), event).Result()
	return
}

func SendEventToUsers(ctx *core.Context, event []byte, users []*proto.User) {
	redisConn, err := ctx.Get("redis.conn")
	if err == nil {
		for _, user := range users {
			_, err := redisConn.(*redis.Client).Publish(ctx, fmt.Sprintf("user:events:%s", user.Id), event).Result()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func SendEventToTheaterMembers(ctx *core.Context, event []byte, theater *models.Theater) (err error) {
	redisConn, err := ctx.Get("redis.conn")
	if err != nil {
		return err
	}
	_, err = redisConn.(*redis.Client).Publish(ctx, fmt.Sprintf("theater:events:%s", theater.ID.Hex()), event).Result()
	return
}

func NewProtoUser(u *models.User) *proto.User {
	return &proto.User{
		Id:            u.ID.Hex(),
		Fullname:      u.Fullname,
		Username:      u.Username,
		Hash:          u.Hash,
		Email:         u.Email,
		IsActive:      u.IsActive,
		IsStaff:       u.IsStaff,
		Verified:      u.Verified,
		EmailVerified: u.EmailVerified,
		Avatar:        u.Avatar,
		TwoFaEnabled:  u.TwoFactorAuthEnabled,
		LastLogin:     timestamppb.New(u.LastLogin),
		JoinedAt:      timestamppb.New(u.JoinedAt),
		UpdatedAt:     timestamppb.New(u.UpdatedAt),
	}
}

func NewProtoUserWithState(user *models.User) *proto.User {
	protoUser := NewProtoUser(user)
	protoUser.State = user.State
	return protoUser
}
