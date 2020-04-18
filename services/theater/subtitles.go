package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func SubtitleToProto(subtitle *models.Subtitle) (*proto.Subtitle, error) {
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

// Get all subtitles from theater
func (s *Service) GetSubtitles(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.TheaterSubtitlesResponse, error) {

	var (
		theater        = new(models.Theater)
		subtitles      = make([]*proto.Subtitle, 0)
		collection     = db.Connection.Collection("subtitles")
		failedResponse = &proto.TheaterSubtitlesResponse{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Could not get subtitles, Please try again later!",
		}
	)

	if _, err := auth.Authenticate(req.AuthRequest); err != nil {
		return &proto.TheaterSubtitlesResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	var (
		theaterObjectID, _ = primitive.ObjectIDFromHex(req.Theater.Id)
		findFilter = bson.M{ "_id": theaterObjectID }
	)

	if err := db.Connection.Collection("theaters").FindOne(ctx, findFilter).Decode(theater); err != nil {
		return &proto.TheaterSubtitlesResponse{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Could not find theater!",
		}, nil
	}

	cursor, err := collection.Find(ctx, bson.M{"theater_id": theaterObjectID})
	if err != nil {
		return failedResponse, nil
	}

	for cursor.Next(ctx) {
		subtitle := new(models.Subtitle)
		if err := cursor.Decode(subtitle); err != nil {
			continue
		}
		protoMsg, err := SubtitleToProto(subtitle)
		if err != nil {
			continue
		}
		subtitles = append(subtitles, protoMsg)
	}

	return &proto.TheaterSubtitlesResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  subtitles,
	}, nil
}

// Remove subtitle from theater
func (s *Service) RemoveSubtitle(ctx context.Context, req *proto.RemoveSubtitleRequest) (*proto.Response, error) {

	var (
		theater        = new(models.Theater)
		collection     = db.Connection.Collection("subtitles")
		failedResponse = &proto.Response{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Could not remove subtitle, Please try again later!",
		}
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	var (
		theaterObjectID, _ = primitive.ObjectIDFromHex(req.Subtitle.TheaterId)
		findFilter = bson.M{
			"_id": theaterObjectID,
			"user_id": user.ID,
		}
	)

	if err := db.Connection.Collection("theaters").FindOne(ctx, findFilter).Decode(theater); err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Could not find theater!",
		}, nil
	}

	var (
		subtitleObjectID, _ = primitive.ObjectIDFromHex(req.Subtitle.Id)
		filter = bson.M{
			"_id": subtitleObjectID,
			"theater_id": theaterObjectID,
		}
	)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount != 1 {
		return failedResponse, err
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Subtitle deleted successfully!",
	}, nil
}
