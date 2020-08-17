package theater

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

func (s *Service) Invite(ctx context.Context, req *proto.InviteFriendsTheaterRequest) (*proto.Response, error) {

	var (
		theater           = new(models.Theater)
		database          = db.Connection
		friends           = make([]*models.User, 0)
		collection        = database.Collection("theaters")
		usersCollection   = database.Collection("users")
		notifsCollections = database.Collection("notifications")
		emptyResponse     = status.Error(codes.Internal, "Could not send invitations, Please tray again later!")
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
		return nil, emptyResponse
	}

	if err := collection.FindOne(ctx, bson.M{ "_id": theaterID }).Decode(&theater); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater!")
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

	cursor, err := usersCollection.Find(ctx, bson.M{"_id": bson.M{"$in": fids}})
	if err != nil {
		return nil, emptyResponse
	}

	for cursor.Next(ctx) {
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
			"type":         int32(proto.Notification_NEW_THEATER_INVITE),
			"read":         false,
			"from_user_id": user.ID,
			"to_user_id":   friend.ID,
			"extra":        theater.ID,
			"read_at":      time.Now(),
			"created_at":   time.Now(),
			"updated_at":   time.Now(),
		})
	}

	if _, err := notifsCollections.InsertMany(ctx, notifications); err != nil {
		return nil, emptyResponse
	}

	for _, friend := range friends {
		// send a new notification event to friend
		err := internal.Client.UserService.SendNewNotificationsEvent(req.AuthRequest, friend.ID.Hex())
		if err != nil {
			sentry.CaptureException(err)
		}
	}

	return &proto.Response{
		Code:     http.StatusOK,
		Status:   "success",
		Message:  "Invitations sent successfully!",
	}, nil
}