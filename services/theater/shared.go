package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func (*Service) GetUserSharedTheaters(ctx context.Context, req *proto.GetAllUserTheatersRequest) (*proto.UserTheatersResponse, error) {

	authUser, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.UserTheatersResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	filter := bson.M{
		"to_user_id": authUser.ID,
		"type": int32(proto.Notification_NEW_THEATER_INVITE),
	}

	qOpts := options.Find()
	qOpts.SetSort(bson.D{
		{"created_at", -1},
	})

	cursor, err := db.Connection.Collection("notifications").Find(mCtx, filter, qOpts)
	if err != nil {
		sentry.CaptureException(err)
		return &proto.UserTheatersResponse{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Could not get any theaters, please try again later!",
		}, nil
	}

	theaterIDs := make([]*primitive.ObjectID, 0)

	for cursor.Next(mCtx) {
		notification := new(models.Notification)
		if err := cursor.Decode(&notification); err != nil {
			continue
		}
		theaterIDs = append(theaterIDs, notification.Extra)
	}

	findTheaters := bson.M{
		"_id": bson.M{
			"$in": theaterIDs,
		},
	}

	thCursor, err := db.Connection.Collection("theaters").Find(mCtx, findTheaters)
	if err != nil {
		sentry.CaptureException(err)
		return &proto.UserTheatersResponse{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Could not get any theaters, please try again later!",
		}, nil
	}

	theaters := make([]*proto.Theater, 0)

	for thCursor.Next(mCtx) {
		theater := new(models.Theater)
		if err := cursor.Decode(theater); err != nil {
			sentry.CaptureException(err)
			break
		}
		th, err := SetDbTheaterToMessageTheater(mCtx, theater)
		if err != nil {
			continue
		}
		theaters = append(theaters, th)
	}

	return &proto.UserTheatersResponse{
		Result:  theaters,
		Code:    http.StatusOK,
		Message: "success",
	}, nil
}