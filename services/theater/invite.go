package theater

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

func (s *Service) Invite(ctx context.Context, req *proto.InviteFriendsTheaterRequest) (*proto.Response, error) {

	var (
		database = db.Connection
		friends     = make([]*models.User, 0)

		collection  = database.Collection("theaters")
		usersCollection = database.Collection("users")
		notificationsCollections = database.Collection("notifications")

		emptyResponse   = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not send invitations, Please tray again later!",
		}

		theater = new(models.Theater)
		mCtx, _ = context.WithTimeout(ctx, 10 * time.Second)
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	theaterID, err := primitive.ObjectIDFromHex(req.TheaterId)
	if err != nil {
		return emptyResponse, err
	}

	if err := collection.FindOne(mCtx, bson.M{ "_id": theaterID }).Decode(&theater); err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Could not find theater!",
		}, err
	}

	fids := make([]primitive.ObjectID, 0)
	for _, friendID := range req.FriendIds {
		if theater.UserId.Hex() == friendID {
			continue
		}
		friendObjectId, err := primitive.ObjectIDFromHex(friendID)
		if err != nil {
			continue
		}
		fids = append(fids, friendObjectId)
	}

	if len(fids) == 0 {
		return &proto.Response{
			Code:     http.StatusOK,
			Status:   "success",
			Message:  "Invitations sent successfully!",
		}, nil
	}

	cursor, err := usersCollection.Find(mCtx, bson.M{"_id": bson.M{"$in": fids}})
	if err != nil {
		return emptyResponse, err
	}

	for cursor.Next(mCtx) {
		var user = new(models.User)
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		if user != nil {
			friends = append(friends, user)
		}
	}

	notifications := make([]interface{}, 0)
	for _, friend := range friends {
		notifications = append(notifications, bson.M{
			"type":         int64(proto.Notification_NEW_THEATER_INVITE),
			"read":         false,
			"from_user_id": user.ID,
			"to_user_id":   friend.ID,
			"extra":        theater.ID,
			"read_at":      time.Now(),
			"created_at":   time.Now(),
			"updated_at":   time.Now(),
		})
	}

	if _, err := notificationsCollections.InsertMany(mCtx, notifications); err != nil {
		return emptyResponse, nil
	}

	return &proto.Response{
		Code:     http.StatusOK,
		Status:   "success",
		Message:  "Invitations sent successfully!",
	}, nil
}