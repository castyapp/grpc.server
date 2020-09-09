package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func (s *Service) UpdateTheater(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.Response, error) {

	var (
		database       = db.Connection
		theater        = new(models.Theater)
		updateMap      = bson.M{}
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

	if req.Theater.Description != "" {
		updateMap["description"] = req.Theater.Description
	}

	if req.Theater.Privacy != proto.PRIVACY_UNKNOWN {
		updateMap["privacy"] = req.Theater.Privacy
	}

	if req.Theater.VideoPlayerAccess != proto.VIDEO_PLAYER_ACCESS_ACCESS_UNKNOWN {
		updateMap["video_player_access"] = req.Theater.VideoPlayerAccess
	}

	if len(updateMap) > 0 {
		updateMap["updated_at"] = time.Now()
		if err := collection.FindOne(ctx, bson.M{ "user_id": user.ID }).Decode(theater); err != nil {
			return nil, failedResponse
		}
		if _, err = collection.UpdateOne(ctx, bson.M{ "_id": theater.ID }, bson.M{ "$set": updateMap }); err != nil {
			return nil, failedResponse
		}
		// sending updated entity through websocket
		//err := internal.Client.TheaterService.SendTheaterUpdateEvent(req.AuthRequest, theater.ID.Hex())
		//if err != nil {
		//	sentry.CaptureException(err)
		//}
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Theater updated successfully!",
	}, nil
}