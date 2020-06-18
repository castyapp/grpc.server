package helpers

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
)

func NewTheaterProto(ctx context.Context, theater *models.Theater) (*proto.Theater, error) {

	var (
		err error
		database           = db.Connection
		thUser             = new(models.User)
		mediaSource        = new(models.MediaSource)
	)

	// finding current media source
	msResult := database.Collection("media_sources").FindOne(ctx, bson.M{"_id": theater.MediaSourceId})
	if err := msResult.Decode(mediaSource); err != nil {
		return nil, err
	}

	// finding theater's user
	uResult := db.Connection.Collection("users").FindOne(ctx, bson.M{ "_id": theater.UserId })
	if err := uResult.Decode(thUser); err != nil {
		return nil, err
	}

	thProtoMessageUser, err := NewProtoUser(thUser)
	if err != nil {
		return nil, err
	}

	mediaSourceProtoMessage, err := NewMediaSourceProto(mediaSource)
	if err != nil {
		return nil, err
	}

	return &proto.Theater{
		Id:                theater.ID.Hex(),
		Title:             theater.Title,
		User:              thProtoMessageUser,
		MediaSource:       mediaSourceProtoMessage,
		Privacy:           theater.Privacy,
		VideoPlayerAccess: theater.VideoPlayerAccess,
	}, nil
}

func NewMediaSourceProto(ms *models.MediaSource) (*proto.MediaSource, error) {
	createdAt, err := ptypes.TimestampProto(ms.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := ptypes.TimestampProto(ms.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &proto.MediaSource{
		Id:               ms.ID.Hex(),
		Type:             ms.Type,
		Banner:           ms.Banner,
		Uri:              ms.Uri,
		LastPlayedTime:   ms.LastPlayedTime,
		Subtitles:        make([]*proto.Subtitle, 0),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}, nil
}

func NewSubtitleProto(subtitle *models.Subtitle) (*proto.Subtitle, error) {
	createdAt, err := ptypes.TimestampProto(subtitle.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := ptypes.TimestampProto(subtitle.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &proto.Subtitle{
		Id: subtitle.ID.Hex(),
		Lang: subtitle.Lang,
		MediaSourceId: subtitle.MediaSourceId.Hex(),
		File: subtitle.File,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}