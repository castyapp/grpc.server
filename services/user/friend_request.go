package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (s *Service) AcceptFriendRequest(ctx context.Context, req *proto.FriendRequest) (*proto.Response, error) {

	var (
		database   = db.Connection
		mCtx, _    = context.WithTimeout(ctx, 20 * time.Second)

		friendRequest = new(models.Friend)
		friendsCollection = database.Collection("friends")

		failedResponse = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not accept friend request, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	frObjectID, err := primitive.ObjectIDFromHex(req.RequestId)
	if err != nil {
		return failedResponse, nil
	}

	filter := bson.M{
		"_id": frObjectID,
		"$or": []interface{}{
			bson.M{"friend_id": user.ID},
			bson.M{"user_id": user.ID},
		},
	}

	if err := friendsCollection.FindOne(mCtx, filter).Decode(&friendRequest); err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Could not find friend request!",
		}, nil
	}

	if friendRequest.Accepted {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Friend request is not valid anymore!",
		}, nil
	}

	update := bson.M{
		"$set": bson.M{
			"accepted": true,
		},
	}
	updateFilter := bson.M{
		"_id": friendRequest.ID,
	}

	updateResult, err := friendsCollection.UpdateOne(mCtx, updateFilter, update)
	if err != nil {
		return failedResponse, nil
	}

	if updateResult.ModifiedCount == 1 {
		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Friend request accepted successfully!",
		}, nil
	}

	return failedResponse, nil
}

func (s *Service) SendFriendRequest(ctx context.Context, req *proto.FriendRequest) (*proto.Response, error) {

	var (
		database   = db.Connection
		friend     = new(models.User)
		mCtx, _    = context.WithTimeout(ctx, 20 * time.Second)

		userCollection    = database.Collection("users")
		friendsCollection = database.Collection("friends")
		notificationsCollection = database.Collection("notifications")

		failedResponse = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not create friend request, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	friendObjectId, err := primitive.ObjectIDFromHex(req.FriendId)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Invalid friend id!",
		}, nil
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

	alreadyFriendsCount, err := friendsCollection.CountDocuments(mCtx, filterFr)
	if err != nil {
		return nil, err
	}

	if alreadyFriendsCount != 0 {
		return &proto.Response{
			Status:  "failed",
			Code:    409,
			Message: "Friend request sent already!",
		}, nil
	}

	if err := userCollection.FindOne(mCtx, bson.M{"_id": friendObjectId}).Decode(friend); err != nil {
		return failedResponse, nil
	}

	friendRequest := bson.M{
		"friend_id": friend.ID,
		"user_id":   user.ID,
		"accepted":  false,
	}

	friendrequestInsertData, err := friendsCollection.InsertOne(mCtx, friendRequest)
	if err != nil {
		return failedResponse, nil
	}

	frInsertID := friendrequestInsertData.InsertedID.(primitive.ObjectID)

	notification := bson.M{
		"type":         int64(proto.NOTIFICATION_TYPE_NEW_FRIEND),
		"read":         false,
		"from_user_id": user.ID,
		"to_user_id":   friend.ID,
		"extra":        frInsertID,
		"read_at":      time.Now(),
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}

	if _, err := notificationsCollection.InsertOne(mCtx, notification); err != nil {
		return failedResponse, nil
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Friend request added successfully!",
	}, nil
}
