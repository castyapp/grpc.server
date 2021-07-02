package theater

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/libcasty-protocol-go/protocol"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Invite(ctx context.Context, req *proto.InviteFriendsTheaterRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                = dbConn.(*mongo.Database)
		theater           = new(models.Theater)
		friends           = make([]*models.User, 0)
		collection        = db.Collection("theaters")
		usersCollection   = db.Collection("users")
		notifsCollections = db.Collection("notifications")
		emptyResponse     = status.Error(codes.Internal, "Could not send invitations, Please tray again later!")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
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

	if err := collection.FindOne(ctx, bson.M{"_id": theaterID}).Decode(&theater); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater!")
	}

	fids := make([]primitive.ObjectID, 0)
	for _, friendID := range req.FriendIds {
		if theater.UserID.Hex() == friendID {
			continue
		}
		friendObjectID, err := primitive.ObjectIDFromHex(friendID)
		if err != nil {
			continue
		}
		fids = append(fids, friendObjectID)
	}

	if len(fids) == 0 {
		return &proto.Response{
			Code:    http.StatusOK,
			Status:  "success",
			Message: "Invitations sent successfully!",
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
		event, err := protocol.NewMsgProtobuf(proto.EMSG_NEW_NOTIFICATION, &proto.NotificationMsgEvent{})
		if err == nil {
			err := helpers.SendEventToUser(s.Context, event.Bytes(), &proto.User{Id: friend.ID.Hex()})
			if err != nil {
				log.Println(err)
			}
		}
	}

	return &proto.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Invitations sent successfully!",
	}, nil
}
