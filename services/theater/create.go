package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (s *Service) CreateTheater(ctx context.Context, req *proto.CreateTheaterRequest) (*proto.Response, error) {

	var (
		collection     = db.Connection.Collection("theaters")
		failedResponse = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not create theater, Please try again later!",
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

	if req.Theater == nil {
		return &proto.Response{
			Status:  "failed",
			Code:    420,
			Message: "Validation error, Theater entry not exists!",
		}, nil
	}

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

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

		subtitles := make([]models.Subtitle, 0)

		theater["movie"] = bson.M{
			"type":             req.Theater.Movie.Type,
			"uri":              movieURI,
			"poster":           req.Theater.Movie.Poster,
			"subtitles":        subtitles,
			"size":             size,
			"length":           length,
			"last_played_time": 0,
		}
	}

	result, err := collection.InsertOne(mCtx, theater)
	if err != nil {
		return failedResponse, nil
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Theater created successfully!",
		Result: []byte(insertedID.Hex()),
	}, nil
}
