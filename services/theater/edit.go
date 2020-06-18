package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/golang/protobuf/ptypes/any"
	"go.mongodb.org/mongo-driver/bson"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func (s *Service) UpdateTheater(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.Response, error) {

	var (
		validationErrors []*any.Any
		database       = db.Connection
		theater        = new(models.Theater)
		collection     = database.Collection("theaters")
		failedResponse = status.Error(codes.Internal, "Could not create theater, Please try again later!")
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

	if req.Theater.MediaSource == nil || req.Theater.MediaSource.Uri == "" {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "movie_uri",
			Value: []byte("MediaSource is required!"),
		})
	}

	if len(validationErrors) > 0 {
		return nil, status.ErrorProto(&spb.Status{
			Code: int32(codes.InvalidArgument),
			Message: "Validation Error!",
			Details: validationErrors,
		})
	}

	if err := collection.FindOne(ctx, bson.M{ "user_id": user.ID }).Decode(theater); err != nil {
		return nil, failedResponse
	}

	_, err = collection.UpdateOne(ctx, bson.M{ "_id": theater.ID }, bson.M{
		"$set": bson.M{
			"title":               req.Theater.Title,
			"privacy":             req.Theater.Privacy,
			"updated_at":          time.Now(),
			"video_player_access": req.Theater.VideoPlayerAccess,
		},
	})

	if err != nil {
		return nil, failedResponse
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Theater updated successfully!",
	}, nil
}