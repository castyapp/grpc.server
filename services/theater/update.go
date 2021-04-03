package theater

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/libcasty-protocol-go/protocol"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpdateTheater(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db             = dbConn.(*mongo.Database)
		theater        = new(models.Theater)
		updateMap      = bson.M{}
		collection     = db.Collection("theaters")
		failedResponse = status.Error(codes.Internal, "Could not create theater, Please try again later!")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
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
		if err := collection.FindOne(ctx, bson.M{"user_id": user.ID}).Decode(theater); err != nil {
			return nil, failedResponse
		}
		if _, err = collection.UpdateOne(ctx, bson.M{"_id": theater.ID}, bson.M{"$set": updateMap}); err != nil {
			return nil, failedResponse
		}
		event, err := protocol.NewMsgProtobuf(proto.EMSG_THEATER_UPDATED, req.Theater)
		if err == nil {
			if err := helpers.SendEventToTheaterMembers(s.Context, event.Bytes(), theater); err != nil {
				log.Println(err)
			}
		}
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Theater updated successfully!",
	}, nil
}
