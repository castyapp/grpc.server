package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (*Service) GetUserSharedTheaters(ctx context.Context, req *proto.GetAllUserTheatersRequest) (*proto.UserTheatersResponse, error) {

	failedErr := status.Error(codes.Internal, "Could not get shared theaters, please try again later!")

	authUser, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"to_user_id": authUser.ID,
		"type": int32(proto.Notification_NEW_THEATER_INVITE),
	}

	qOpts := options.Find()
	qOpts.SetSort(bson.D{
		{"created_at", -1},
	})

	cursor, err := db.Connection.Collection("notifications").Find(ctx, filter, qOpts)
	if err != nil {
		return nil, failedErr
	}

	theaterIDs := make([]*primitive.ObjectID, 0)

	for cursor.Next(ctx) {
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

	thCursor, err := db.Connection.Collection("theaters").Find(ctx, findTheaters)
	if err != nil {
		return nil, failedErr
	}

	theaters := make([]*proto.Theater, 0)

	for thCursor.Next(ctx) {
		theater := new(models.Theater)
		if err := thCursor.Decode(theater); err != nil {
			continue
		}
		th, err := helpers.NewTheaterProto(ctx, theater)
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