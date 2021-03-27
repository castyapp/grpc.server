package theater

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get all subtitles from theater
func (s *Service) AddSubtitles(ctx context.Context, req *proto.AddSubtitlesRequest) (*proto.SubtitlesResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                     = dbConn.(*mongo.Database)
		insertMap              = make([]interface{}, 0)
		mediaSource            = new(models.MediaSource)
		mediaSourcesCollection = db.Collection("media_sources")
		subtitlesCollection    = db.Collection("subtitles")
		failedResponse         = status.Error(codes.Internal, "Could not add subtitles, Please try again later!")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	var (
		mediaSourceObjectID, _ = primitive.ObjectIDFromHex(req.MediaSourceId)
		findFilter             = bson.M{
			"_id":     mediaSourceObjectID,
			"user_id": user.ID,
		}
	)

	if err := mediaSourcesCollection.FindOne(ctx, findFilter).Decode(mediaSource); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find media source!")
	}

	for _, subtitle := range req.Subtitles {
		insertMap = append(insertMap, bson.M{
			"user_id":         user.ID,
			"media_source_id": mediaSource.ID,
			"file":            subtitle.File,
			"lang":            subtitle.Lang,
			"created_at":      time.Now(),
			"updated_at":      time.Now(),
		})
	}

	if _, err := subtitlesCollection.InsertMany(ctx, insertMap); err != nil {
		log.Println(err)
		return nil, failedResponse
	}

	return &proto.SubtitlesResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Subtitle added successfully!",
	}, nil
}

// Get all subtitles from theater
func (s *Service) GetSubtitles(ctx context.Context, req *proto.MediaSourceAuthRequest) (*proto.TheaterSubtitlesResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                     = dbConn.(*mongo.Database)
		authenticated          = false
		authUser               = new(models.User)
		dbTheater              = new(models.Theater)
		mediaSource            = new(models.MediaSource)
		subtitles              = make([]*proto.Subtitle, 0)
		theatersCollection     = db.Collection("theaters")
		mediaSourcesCollection = db.Collection("media_sources")
		collection             = db.Collection("subtitles")
		failedResponse         = status.Error(codes.Internal, "Could not get subtitles, Please try again later!")
	)

	if req.AuthRequest != nil {
		if authUser, err = auth.Authenticate(s.Context, req.AuthRequest); err != nil {
			return nil, err
		}
		authenticated = true
	}

	mediaSourceObjectID, err := primitive.ObjectIDFromHex(req.Media.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not parse theater id!")
	}

	if err := theatersCollection.FindOne(ctx, bson.M{"media_source_id": mediaSourceObjectID}).Decode(dbTheater); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater with this media source!")
	}

	if !authenticated {
		switch dbTheater.Privacy {
		case proto.PRIVACY_PRIVATE:
			return nil, status.Error(codes.PermissionDenied, "Permission Denied!")
		}
	} else {
		if dbTheater.UserId.Hex() != authUser.ID.Hex() {
			switch dbTheater.Privacy {
			case proto.PRIVACY_PRIVATE:
				return nil, status.Error(codes.PermissionDenied, "Permission Denied!")
			}
		}
	}

	if err := mediaSourcesCollection.FindOne(ctx, bson.M{"_id": dbTheater.MediaSourceId}).Decode(mediaSource); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find media source!")
	}

	cursor, err := collection.Find(ctx, bson.M{"media_source_id": mediaSource.ID})
	if err != nil {
		return nil, failedResponse
	}

	for cursor.Next(ctx) {
		subtitle := new(models.Subtitle)
		if err := cursor.Decode(subtitle); err != nil {
			continue
		}
		protoMsg, err := helpers.NewSubtitleProto(subtitle)
		if err != nil {
			continue
		}
		subtitles = append(subtitles, protoMsg)
	}

	return &proto.TheaterSubtitlesResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: subtitles,
	}, nil
}

// Remove subtitle from theater
func (s *Service) RemoveSubtitle(ctx context.Context, req *proto.RemoveSubtitleRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db             = dbConn.(*mongo.Database)
		mediaSource    = new(models.Theater)
		collection     = db.Collection("subtitles")
		failedResponse = status.Error(codes.Internal, "Could not remove subtitle, Please try again later!")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	var (
		mediaSourceObjectID, _ = primitive.ObjectIDFromHex(req.MediaSourceId)
		findFilter             = bson.M{
			"_id":     mediaSourceObjectID,
			"user_id": user.ID,
		}
	)

	if err := db.Collection("media_sources").FindOne(ctx, findFilter).Decode(mediaSource); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find media source!")
	}

	var (
		subtitleObjectID, _ = primitive.ObjectIDFromHex(req.SubtitleId)
		filter              = bson.M{
			"_id":             subtitleObjectID,
			"media_source_id": mediaSource.ID,
		}
	)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount != 1 {
		return nil, failedResponse
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Subtitle deleted successfully!",
	}, nil
}
