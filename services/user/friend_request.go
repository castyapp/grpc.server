package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/internal"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func (s *Service) AcceptFriendRequest(ctx context.Context, req *proto.FriendRequest) (*proto.Response, error) {

	var (
		friendRequest     = new(models.Friend)
		friendsCollection = db.Connection.Collection("friends")
		notifsCollection  = db.Connection.Collection("notifications")
		failedResponse    = status.Error(codes.Internal, "Could not accept friend request, Please try again later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	frObjectID, err := primitive.ObjectIDFromHex(req.RequestId)
	if err != nil {
		return nil, failedResponse
	}

	filter := bson.M{
		"_id": frObjectID,
		"$or": []interface{}{
			bson.M{"friend_id": user.ID},
			bson.M{"user_id": user.ID},
		},
	}

	if err := friendsCollection.FindOne(ctx, filter).Decode(&friendRequest); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find friend request!")
	}

	if friendRequest.Accepted {
		return nil, status.Error(codes.InvalidArgument, "Friend request is not valid anymore!")
	}

	findNotif := bson.M{
		"extra": friendRequest.ID,
		"to_user_id": user.ID,
	}

	// update user's notification to read
	_, _ = notifsCollection.UpdateOne(ctx, findNotif, bson.M{
		"$set": bson.M{
			"read": true,
			"updated_at": time.Now(),
			"read_at": time.Now(),
		},
	})

	var (
		update = bson.M{
			"$set": bson.M{
				"accepted": true,
			},
		}
		updateFilter = bson.M{
			"_id": friendRequest.ID,
		}
	)

	updateResult, err := friendsCollection.UpdateOne(ctx, updateFilter, update)
	if err != nil {
		return nil, failedResponse
	}

	if updateResult.ModifiedCount == 1 {

		friendID := friendRequest.FriendId.Hex()
		if friendRequest.FriendId.Hex() == user.ID.Hex() {
			friendID = friendRequest.UserId.Hex()
		}

		// send new friend request event to friend websocket clients
		_ = internal.Client.UserService.AcceptNotificationEvent(req.AuthRequest, user, friendID)

		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Friend request accepted successfully!",
		}, nil
	}

	return nil, failedResponse
}

func (s *Service) SendFriendRequest(ctx context.Context, req *proto.FriendRequest) (*proto.Response, error) {

	var (
		database                = db.Connection
		friend                  = new(models.User)
		userCollection          = database.Collection("users")
		friendsCollection       = database.Collection("friends")
		notificationsCollection = database.Collection("notifications")
		failedResponse          = status.Error(codes.Internal, "Could not create friend request, Please try again later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	friendObjectId, err := primitive.ObjectIDFromHex(req.FriendId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid friend id!")
	}

	var (
		filterFr = bson.M{
			"$or": []interface{}{
				bson.M{
					"friend_id": user.ID,
					"user_id": friendObjectId,
				},
				bson.M{
					"user_id": user.ID,
					"friend_id": friendObjectId,
				},
			},
		}
	)

	alreadyFriendsCount, err := friendsCollection.CountDocuments(ctx, filterFr)
	if err != nil {
		return nil, failedResponse
	}

	if alreadyFriendsCount != 0 {
		return nil, status.Error(codes.Aborted, "Friend request sent already!")
	}

	if err := userCollection.FindOne(ctx, bson.M{"_id": friendObjectId}).Decode(friend); err != nil {
		return nil, failedResponse
	}

	friendRequest := bson.M{
		"friend_id": friend.ID,
		"user_id":   user.ID,
		"accepted":  false,
	}

	friendrequestInsertData, err := friendsCollection.InsertOne(ctx, friendRequest)
	if err != nil {
		return nil, failedResponse
	}

	frInsertID := friendrequestInsertData.InsertedID.(primitive.ObjectID)

	notification := bson.M{
		"type":         int64(proto.Notification_NEW_FRIEND),
		"read":         false,
		"from_user_id": user.ID,
		"to_user_id":   friend.ID,
		"extra":        frInsertID,
		"read_at":      time.Now(),
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}

	if _, err := notificationsCollection.InsertOne(ctx, notification); err != nil {
		return nil, failedResponse
	}

	// send new friend request event to friend websocket clients
	err = internal.Client.UserService.SendNewNotificationsEvent(friend.ID.Hex())
	if err != nil {
		sentry.CaptureException(err)
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Friend request added successfully!",
	}, nil
}
