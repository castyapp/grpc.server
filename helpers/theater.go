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

		thUser  = new(models.User)
		thProtoMessageUser = new(proto.User)

		movie = &proto.Movie{
			Type:            theater.Movie.Type,
			Uri:             theater.Movie.Uri,
			Poster:          theater.Movie.Poster,
			Size:            int64(theater.Movie.Size),
			Length:          int64(theater.Movie.Length),
			LastPlayedTime:  theater.Movie.LastPlayedTime,
			Subtitles:       []*proto.Subtitle{},
		}
	)

	find := db.Connection.Collection("users").FindOne(ctx, bson.M{ "_id": theater.UserId })
	if err := find.Decode(thUser); err != nil {
		return nil, err
	}

	thProtoMessageUser, err = NewProtoUser(thUser)
	if err != nil {
		return nil, err
	}

	return &proto.Theater{
		Id:      theater.ID.Hex(),
		Title:   theater.Title,
		Hash:    theater.Hash,
		User:    thProtoMessageUser,
		Movie:   movie,
		Privacy: proto.PRIVACY(theater.Privacy),
		VideoPlayerAccess: proto.PRIVACY(theater.VideoPlayerAccess),
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
		TheaterId: subtitle.TheaterId.Hex(),
		File: subtitle.File,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}