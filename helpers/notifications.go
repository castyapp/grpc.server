package helpers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewNotificationProto(db *mongo.Database, notif *models.Notification) (*proto.Notification, error) {

	var (
		readAt, _    = ptypes.TimestampProto(notif.ReadAt)
		createdAt, _ = ptypes.TimestampProto(notif.CreatedAt)
		updatedAt, _ = ptypes.TimestampProto(notif.UpdatedAt)
		fromUser     = new(models.User)
		mCtx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	)

	defer cancel()

	cursor := db.Collection("users").FindOne(mCtx, bson.M{
		"_id": notif.FromUserId,
	})
	if err := cursor.Decode(&fromUser); err != nil {
		return nil, err
	}

	protoMSG := &proto.Notification{
		Id:        notif.ID.Hex(),
		Type:      notif.Type,
		Read:      notif.Read,
		ReadAt:    readAt,
		FromUser:  NewProtoUser(fromUser),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	switch notif.Type {
	case proto.Notification_NEW_FRIEND:
		notifFriendData := new(models.Friend)
		cursor := db.Collection("friends").FindOne(mCtx, bson.M{
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
		cursor := db.Collection("theaters").FindOne(mCtx, bson.M{
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
