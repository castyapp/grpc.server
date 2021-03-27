package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/castyapp/grpc.server/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationData struct {
	Data string `json:"data"`
	User string `json:"user"`
}

func (s *Service) CreateNotification(ctx context.Context, req *proto.CreateNotificationRequest) (*proto.NotificationResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db             = dbConn.(*mongo.Database)
		collection     = db.Collection("notifications")
		failedResponse = status.Error(codes.InvalidArgument, "Could not create notification, Please try again later!")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	if req.Notification == nil {
		return nil, status.Error(codes.InvalidArgument, "Notification entry not exists!")
	}

	friendObjectId, err := primitive.ObjectIDFromHex(req.Notification.ToUserId)
	if err != nil {
		return nil, failedResponse
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
			return nil, failedResponse
		}
		theaterObjectId, err := primitive.ObjectIDFromHex(friend.Id)
		if err != nil {
			return nil, failedResponse
		}
		notification["extra"] = theaterObjectId
	}

	result, err := collection.InsertOne(ctx, notification)
	if err != nil {
		return nil, failedResponse
	}

	var (
		insertedID   = result.InsertedID.(*primitive.ObjectID)
		createdAt, _ = ptypes.TimestampProto(time.Now())
	)

	return &proto.NotificationResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Notification created successfully!",
		Result: []*proto.Notification{
			{
				Id:         insertedID.Hex(),
				Type:       req.Notification.Type,
				Data:       notification["extra"].(string),
				Read:       false,
				FromUserId: user.ID.Hex(),
				ToUserId:   friendObjectId.Hex(),
				CreatedAt:  createdAt,
			},
		},
	}, nil
}

func (s *Service) GetNotifications(ctx context.Context, req *proto.AuthenticateRequest) (*proto.NotificationResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db               = dbConn.(*mongo.Database)
		notifications    = make([]*proto.Notification, 0)
		notifsCollection = db.Collection("notifications")
		failedResponse   = status.Error(codes.Internal, "Could not get notifications, Please try again later!")
	)

	user, err := auth.Authenticate(s.Context, req)
	if err != nil {
		return nil, err
	}

	qOpts := options.Find()
	qOpts.SetSort(bson.D{
		{"created_at", -1},
	})

	cursor, err := notifsCollection.Find(ctx, bson.M{"to_user_id": user.ID}, qOpts)
	if err != nil {
		return nil, failedResponse
	}

	var unreadCount int64 = 0

	for cursor.Next(ctx) {
		notification := new(models.Notification)
		if err := cursor.Decode(notification); err != nil {
			continue
		}
		messageNotification, err := helpers.NewNotificationProto(db, notification)
		if err != nil {
			continue
		}
		if !notification.Read {
			unreadCount++
		}
		notifications = append(notifications, messageNotification)
	}

	return &proto.NotificationResponse{
		Status:      "success",
		Code:        http.StatusOK,
		Result:      notifications,
		UnreadCount: unreadCount,
	}, nil
}

func (s *Service) ReadAllNotifications(ctx context.Context, req *proto.AuthenticateRequest) (*proto.NotificationResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db               = dbConn.(*mongo.Database)
		notifsCollection = db.Collection("notifications")
		failedResponse   = status.Error(codes.Internal, "Could not update notifications, Please try again later!")
	)

	user, err := auth.Authenticate(s.Context, req)
	if err != nil {
		return nil, err
	}

	var (
		filter = bson.M{
			"to_user_id": user.ID,
			"read":       false,
		}
		update = bson.M{
			"$set": bson.M{
				"read": true,
			},
		}
	)

	if _, err := notifsCollection.UpdateMany(ctx, filter, update); err != nil {
		return nil, failedResponse
	}

	return &proto.NotificationResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Notifications updated successfully!",
	}, nil
}
