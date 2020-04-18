package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/services"
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

func (s *Service) CreateTheater(ctx context.Context, req *proto.CreateTheaterRequest) (*proto.Response, error) {

	var (
		collection     = db.Connection.Collection("theaters")
		failedResponse = status.Error(codes.Internal, "Could not create theater, Please try again later!")
		validationErrors []*any.Any
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized!")
	}

	if req.Theater == nil {
		return nil, status.Error(codes.InvalidArgument, "Validation error, Theater entry not exists!")
	}

	if req.Theater.Title == "" {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "title",
			Value: []byte("Title is required!"),
		})
	}

	if req.Theater.Movie == nil || req.Theater.Movie.Uri == "" {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "movie_uri",
			Value: []byte("MovieUri is required!"),
		})
	}

	if len(validationErrors) > 0 {
		return nil, status.ErrorProto(&spb.Status{
			Code: int32(codes.InvalidArgument),
			Message: "Validation Error!",
			Details: validationErrors,
		})
	}

	theater := bson.M{
		"title":      req.Theater.Title,
		"hash":       services.GenerateHash(),
		"privacy":    int64(req.Theater.Privacy),
		"user_id":    user.ID,
		"created_at": time.Now(),
		"updated_at": time.Now(),
		"video_player_access": int64(req.Theater.VideoPlayerAccess),
	}

	if req.Theater.Movie != nil {

		var (
			size int64 = 0
			length int64 = 0
			movieURI = req.Theater.Movie.Uri
		)

		switch movieTYPE := req.Theater.Movie.Type; movieTYPE {
		case proto.MovieType_URI:
			movieDuration, err := GetMovieDuration(movieURI)
			if err == nil {
				length = movieDuration
			}
			movieSize, err := GetMovieFileSize(movieURI)
			if err == nil {
				size = movieSize
			}
		}

		theater["movie"] = bson.M{
			"type":             req.Theater.Movie.Type,
			"uri":              movieURI,
			"poster":           "default",
			"size":             size,
			"length":           length,
			"last_played_time": 0,
		}
	}

	result, err := collection.InsertOne(ctx, theater)
	if err != nil {
		return nil, failedResponse
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Theater created successfully!",
		Result: []byte(insertedID.Hex()),
	}, nil
}
