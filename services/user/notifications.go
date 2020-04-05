package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetDBNotificationToProto(notif *models.Notification) (*proto.Notification, error) {

	var (
		readAt, _    = ptypes.TimestampProto(notif.ReadAt)
		createdAt, _ = ptypes.TimestampProto(notif.CreatedAt)
		updatedAt, _ = ptypes.TimestampProto(notif.UpdatedAt)
		fromUser     = new(models.User)
		mCtx, _      = context.WithTimeout(context.Background(), 10*time.Second)
	)

	cursor := db.Connection.Collection("users").FindOne(mCtx, bson.M{
		"_id": notif.FromUserId,
	})
	if err := cursor.Decode(&fromUser); err != nil {
		return nil, err
	}

	protoUser, err := SetDBUserToProtoUser(fromUser)
	if err != nil {
		return nil, err
	}

	protoMSG := &proto.Notification{
		Id:        notif.ID.Hex(),
		Type:      notif.Type,
		Read:      notif.Read,
		ReadAt:    readAt,
		FromUser:  protoUser,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	switch notif.Type {
	case proto.Notification_NEW_FRIEND:
		notifFriendData := new(models.Friend)
		cursor := db.Connection.Collection("friends").FindOne(mCtx, bson.M{
			"_id": notif.Extra,
		})
		if err := cursor.Decode(&notifFriendData); err != nil {
			return nil, err
		}
		ntfJson, err := json.Marshal(notifFriendData)
		if err != nil {
			return nil, err
		}
		protoMSG.Data = string(ntfJson)
	case proto.Notification_NEW_THEATER_INVITE:
		notifTheaterData := new(models.Theater)
		cursor := db.Connection.Collection("theaters").FindOne(mCtx, bson.M{
			"_id": notif.Extra,
		})
		if err := cursor.Decode(&notifTheaterData); err != nil {
			return nil, err
		}
		ntfJson, err := json.Marshal(notifTheaterData)
		if err != nil {
			return nil, err
		}
		protoMSG.Data = string(ntfJson)
	}

	return protoMSG, nil
}

type NotificationData struct {
	Data string `json:"data"`
	User string `json:"user"`
}

func (s *Service) CreateNotification(ctx context.Context, req *proto.CreateNotificationRequest) (*proto.NotificationResponse, error) {

	var (
		mCtx, _        = context.WithTimeout(ctx, 20*time.Second)
		collection     = db.Connection.Collection("notifications")
		failedResponse = &proto.NotificationResponse{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not create notification, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.NotificationResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	if req.Notification == nil {
		return &proto.NotificationResponse{
			Status:  "failed",
			Code:    420,
			Message: "Validation error, Notification entry not exists!",
		}, nil
	}

	friendObjectId, err := primitive.ObjectIDFromHex(req.Notification.ToUserId)
	if err != nil {
		return failedResponse, nil
	}

	notification := bson.M{
		"type":         int64(req.Notification.Type),
		"read":         req.Notification.Read,
		"from_user_id": user.ID,
		"to_user_id":   friendObjectId,
		"read_at":      time.Now(),
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}

	switch req.Notification.Type {
	case proto.Notification_NEW_THEATER_INVITE:

		friend := new(proto.User)
		err := json.Unmarshal([]byte(req.Notification.Data), friend)
		if err != nil {
			return failedResponse, nil
		}

		theaterObjectId, err := primitive.ObjectIDFromHex(friend.Id)
		if err != nil {
			return failedResponse, nil
		}

		notification["extra"] = theaterObjectId
	}

	if _, err := collection.InsertOne(mCtx, notification); err != nil {
		return failedResponse, nil
	}

	return &proto.NotificationResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Notification created successfully!",
	}, nil
}

func (s *Service) GetNotifications(ctx context.Context, req *proto.AuthenticateRequest) (*proto.NotificationResponse, error) {

	var (
		notifications []*proto.Notification

		database = db.Connection
		mCtx, _  = context.WithTimeout(ctx, 20*time.Second)

		notificationCollection = database.Collection("notifications")

		failedResponse = &proto.NotificationResponse{
			Status:      "failed",
			Code:        http.StatusInternalServerError,
			UnreadCount: 0,
			Message:     "Could not get notifications, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req)
	if err != nil {
		return &proto.NotificationResponse{
			Status:      "failed",
			Code:        http.StatusUnauthorized,
			Message:     "Unauthorized!",
			UnreadCount: 0,
		}, nil
	}

	qOpts := options.Find()
	qOpts.SetSort(bson.D{
		{"created_at", -1},
	})

	cursor, err := notificationCollection.Find(mCtx, bson.M{"to_user_id": user.ID}, qOpts)
	if err != nil {
		return failedResponse, nil
	}

	for cursor.Next(mCtx) {

		notification := new(models.Notification)
		if err := cursor.Decode(notification); err != nil {
			break
		}

		messageNotification, err := SetDBNotificationToProto(notification)
		if err != nil {
			break
		}

		notifications = append(notifications, messageNotification)
	}

	filter := bson.M{"to_user_id": user.ID, "read": false}
	unreadCount, err := notificationCollection.CountDocuments(mCtx, filter)
	if err != nil {
		return failedResponse, nil
	}

	return &proto.NotificationResponse{
		Status:      "success",
		Code:        http.StatusOK,
		Result:      notifications,
		UnreadCount: unreadCount,
	}, nil
}

func (s *Service) ReadAllNotifications(ctx context.Context, req *proto.AuthenticateRequest) (*proto.NotificationResponse, error) {

	var (
		mCtx, _                = context.WithTimeout(ctx, 10 * time.Second)
		notificationCollection = db.Connection.Collection("notifications")
		failedResponse         = &proto.NotificationResponse{
			Status:      "failed",
			Code:        http.StatusInternalServerError,
			UnreadCount: 0,
			Message:     "Could not update notifications, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req)
	if err != nil {
		return &proto.NotificationResponse{
			Status:      "failed",
			Code:        http.StatusUnauthorized,
			Message:     "Unauthorized!",
			UnreadCount: 0,
		}, nil
	}

	var (
		filter = bson.M{
			"to_user_id": user.ID,
			"read": false,
		}
		update = bson.M{
			"$set": bson.M{
				"read": true,
			},
		}
	)

	if _, err := notificationCollection.UpdateMany(mCtx, filter, update); err != nil {
		return failedResponse, nil
	}

	return &proto.NotificationResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Notifications updated successfully!",
	}, nil
}
