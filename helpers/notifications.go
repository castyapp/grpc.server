package helpers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewNotificationProto(db *mongo.Database, n *models.Notification) (*proto.Notification, error) {

	var (
		fromUser     = new(models.User)
		mCtx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	)

	defer cancel()

	cursor := db.Collection("users").FindOne(mCtx, bson.M{"_id": n.FromUserID})
	if err := cursor.Decode(&fromUser); err != nil {
		return nil, err
	}

	protoMSG := &proto.Notification{
		Id:        n.ID.Hex(),
		Type:      n.Type,
		Read:      n.Read,
		ReadAt:    timestamppb.New(n.ReadAt),
		CreatedAt: timestamppb.New(n.ReadAt),
		UpdatedAt: timestamppb.New(n.UpdatedAt),
		FromUser:  NewProtoUser(fromUser),
	}

	switch n.Type {
	case proto.Notification_NEW_FRIEND:
		notifFriendData := new(models.Friend)
		cursor := db.Collection("friends").FindOne(mCtx, bson.M{
			"_id": n.Extra,
		})
		if err := cursor.Decode(&notifFriendData); err != nil {
			return nil, err
		}
		ntfJSON, err := json.Marshal(notifFriendData)
		if err != nil {
			return nil, err
		}
		protoMSG.Data = string(ntfJSON)
	case proto.Notification_NEW_THEATER_INVITE:
		notifTheaterData := new(models.Theater)
		cursor := db.Collection("theaters").FindOne(mCtx, bson.M{
			"_id": n.Extra,
		})
		if err := cursor.Decode(&notifTheaterData); err != nil {
			return nil, err
		}
		ntfJSON, err := json.Marshal(notifTheaterData)
		if err != nil {
			return nil, err
		}
		protoMSG.Data = string(ntfJSON)
	}

	return protoMSG, nil
}
