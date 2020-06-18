package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes/any"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func (s *Service) AddMediaSource(ctx context.Context, req *proto.MediaSourceAuthRequest) (*proto.Response, error) {

	var (
		validationErrors []*any.Any
		database   = db.Connection
		collection = database.Collection("media_sources")
		failedResponse = status.Error(codes.Internal, "Could not add a new media source. Please try agian later!")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized!")
	}

	if req.Media.Uri == "" {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "uri",
			Value: []byte("Uri is required!"),
		})
	}

	if req.Media.Type == proto.MediaSource_UNKNOWN {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "type",
			Value: []byte("Media source type can not be unknown!"),
		})
	}

	if len(validationErrors) > 0 {
		return nil, status.ErrorProto(&spb.Status{
			Code: int32(codes.InvalidArgument),
			Message: "Validation Error!",
			Details: validationErrors,
		})
	}

	mediaSource := bson.M{
		"type":               req.Media.Type,
		"banner":             req.Media.Banner,
		"uri":                req.Media.Uri,
		"last_played_time":   req.Media.LastPlayedTime,
		"user_id":            user.ID,
		"created_at":         time.Now(),
		"updated_at":         time.Now(),
	}

	if _, err := collection.InsertOne(ctx, mediaSource); err != nil {
		return nil, failedResponse
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Media source created successfully!",
	}, nil
}

func (s *Service) GetMediaSources(ctx context.Context, req *proto.MediaSourceAuthRequest) (*proto.TheaterMediaSourcesResponse, error) {

	var (
		database     = db.Connection
		mediaSources = make([]*proto.MediaSource, 0)
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.TheaterMediaSourcesResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	cursor, err := database.Collection("media_sources").Find(ctx, bson.M{ "user_id": user.ID })
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater!")
	}

	for cursor.Next(ctx) {
		dbMediaSource := new(models.MediaSource)
		if err := cursor.Decode(dbMediaSource); err != nil {
			continue
		}
		protoMediaSource, err := helpers.NewMediaSourceProto(dbMediaSource)
		if err != nil {
			continue
		}
		mediaSources = append(mediaSources, protoMediaSource)
	}

	return &proto.TheaterMediaSourcesResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  mediaSources,
	}, nil
}

func (s *Service) RemoveMediaSource(ctx context.Context, req *proto.MediaSourceRemoveRequest) (*proto.Response, error) {

	collection := db.Connection.Collection("media_sources")

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	mediaSourceObjectID, err := primitive.ObjectIDFromHex(req.MediaSourceId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not parse MediaSourceId!")
	}

	if result, err := collection.DeleteOne(ctx, bson.M{ "_id": mediaSourceObjectID, "user_id": user.ID }); err == nil {
		if result.DeletedCount == 1 {
			return &proto.Response{
				Status:  "success",
				Code:    http.StatusOK,
				Message: "Media source deleted successfully@",
			}, nil
		}
	}

	return nil, status.Error(codes.Aborted, "Could not delete media source. Please try again later!")
}
