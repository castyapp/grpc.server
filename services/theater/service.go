package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type Service struct {}

func SetDbTheaterToMessageTheater(ctx context.Context, theater *models.Theater) (*proto.Theater, error) {

	var (
		err error

		thUser  = new(models.User)
		thProtoMessageUser = new(proto.User)

		movie = &proto.Movie{
			Type:            theater.Movie.Type,
			Uri:             theater.Movie.Uri,
			Poster:          theater.Movie.Poster,
			Size:            int64(theater.Movie.Size),
			Length:          int64(theater.Movie.Length),
			LastPlayedTime:  theater.Movie.LastPlayedTime,
			Subtitles:       []*proto.Subtitle{},
		}
	)

	find := db.Connection.Collection("users").FindOne(ctx, bson.M{ "_id": theater.UserId })
	if err := find.Decode(thUser); err != nil {
		return nil, err
	}

	thProtoMessageUser, err = helpers.SetDBUserToProtoUser(thUser)
	if err != nil {
		return nil, err
	}

	return &proto.Theater{
		Id:      theater.ID.Hex(),
		Title:   theater.Title,
		Hash:    theater.Hash,
		User:    thProtoMessageUser,
		Movie:   movie,
		Privacy: proto.PRIVACY(theater.Privacy),
		VideoPlayerAccess: proto.PRIVACY(theater.VideoPlayerAccess),
	}, nil
}

func (s *Service) GetUserTheaters(ctx context.Context, req *proto.GetAllUserTheatersRequest) (*proto.UserTheatersResponse, error) {

	var (
		theaters   = make([]*proto.Theater, 0)
		collection = db.Connection.Collection("theaters")
	)

	authUser, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.UserTheatersResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	qOpts := options.Find()
	qOpts.SetSort(bson.D{
		{"created_at", -1},
	})

	cursor, err := collection.Find(ctx, bson.M{"user_id": authUser.ID}, qOpts)
	if err != nil {
		sentry.CaptureException(err)
		return &proto.UserTheatersResponse{
			Status:  "failed",
			Code:    http.StatusNotAcceptable,
			Message: "The requested parameter is not acceptable!",
		}, nil
	}

	for cursor.Next(ctx) {
		theater := new(models.Theater)
		if err := cursor.Decode(theater); err != nil {
			sentry.CaptureException(err)
			break
		}
		th, err := SetDbTheaterToMessageTheater(ctx, theater)
		if err != nil {
			break
		}
		theaters = append(theaters, th)
	}

	return &proto.UserTheatersResponse{
		Result:  theaters,
		Code:    http.StatusOK,
		Message: "success",
	}, nil
}