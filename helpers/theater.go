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
		database                 = db.Connection
		thUser                   = new(models.User)
		mediaSourceProtoMessage  = new(proto.MediaSource)
		mediaSource              = new(models.MediaSource)
	)

	if theater.MediaSourceId != nil {
		// finding current media source
		msResult := database.Collection("media_sources").FindOne(ctx, bson.M{"_id": theater.MediaSourceId})
		if err := msResult.Decode(mediaSource); err == nil {
			mediaSourceProtoMessage = NewMediaSourceProto(mediaSource)
		}
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

	return &proto.Theater{
		Id:                theater.ID.Hex(),
		Description:       theater.Description,
		User:              thProtoMessageUser,
		MediaSource:       mediaSourceProtoMessage,
		Privacy:           theater.Privacy,
		VideoPlayerAccess: theater.VideoPlayerAccess,
	}, nil
}

func NewMediaSourceProto(ms *models.MediaSource) *proto.MediaSource {
	createdAt, _ := ptypes.TimestampProto(ms.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(ms.UpdatedAt)
	return &proto.MediaSource{
		Id:               ms.ID.Hex(),
		Title:            ms.Title,
		Type:             ms.Type,
		Banner:           ms.Banner,
		Uri:              ms.Uri,
		Length:           ms.Length,
		Artist:           ms.Artist,
		Subtitles:        make([]*proto.Subtitle, 0),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
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