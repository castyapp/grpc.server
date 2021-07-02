package helpers

import (
	"context"

	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewTheaterProto(ctx context.Context, db *mongo.Database, theater *models.Theater) (*proto.Theater, error) {

	var (
		thUser                  = new(models.User)
		mediaSourceProtoMessage = new(proto.MediaSource)
		mediaSource             = new(models.MediaSource)
	)

	if theater.MediaSourceID != nil {
		// finding current media source
		msResult := db.Collection("media_sources").FindOne(ctx, bson.M{"_id": theater.MediaSourceID})
		if err := msResult.Decode(mediaSource); err == nil {
			mediaSourceProtoMessage = NewMediaSourceProto(mediaSource)
		}
	}

	// finding theater's user
	uResult := db.Collection("users").FindOne(ctx, bson.M{"_id": theater.UserID})
	if err := uResult.Decode(thUser); err != nil {
		return nil, err
	}

	return &proto.Theater{
		Id:                theater.ID.Hex(),
		Description:       theater.Description,
		User:              NewProtoUser(thUser),
		MediaSource:       mediaSourceProtoMessage,
		Privacy:           theater.Privacy,
		VideoPlayerAccess: theater.VideoPlayerAccess,
	}, nil
}

func NewMediaSourceProto(ms *models.MediaSource) *proto.MediaSource {
	createdAt := timestamppb.New(ms.CreatedAt)
	updatedAt := timestamppb.New(ms.UpdatedAt)
	return &proto.MediaSource{
		Id:        ms.ID.Hex(),
		Title:     ms.Title,
		Type:      ms.Type,
		Banner:    ms.Banner,
		Uri:       ms.URI,
		Length:    ms.Length,
		Artist:    ms.Artist,
		Subtitles: make([]*proto.Subtitle, 0),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewSubtitleProto(s *models.Subtitle) (*proto.Subtitle, error) {
	createdAt := timestamppb.New(s.CreatedAt)
	updatedAt := timestamppb.New(s.UpdatedAt)
	return &proto.Subtitle{
		Id:            s.ID.Hex(),
		Lang:          s.Lang,
		MediaSourceId: s.MediaSourceID.Hex(),
		File:          s.File,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}
