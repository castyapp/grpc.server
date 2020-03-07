package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.proto/messages"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/CastyLab/grpc.server/services/user"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type Service struct {}

func SetDbTheaterToMessageTheater(ctx context.Context, theater *models.Theater) (*messages.Theater, error) {

	var (
		err error

		thUser  = new(models.User)
		thProtoMessageUser = new(messages.User)

		movie = &messages.Movie{
			MovieUri: theater.Movie.MovieUri,
			Poster: theater.Movie.Poster,
			Size:   int64(theater.Movie.Size),
			Length: int64(theater.Movie.Length),
			LastPlayedTime: theater.Movie.LastPlayedTime,
			Subtitles: []*messages.Subtitle{},
		}
	)

	find := db.Connection.Collection("users").FindOne(ctx, bson.M{ "_id": theater.UserId })
	if err := find.Decode(thUser); err != nil {
		return nil, err
	}

	thProtoMessageUser, err = user.SetDBUserToProtoUser(thUser)
	if err != nil {
		return nil, err
	}

	return &messages.Theater{
		Id:      theater.ID.Hex(),
		Title:   theater.Title,
		Hash:    theater.Hash,
		User:    thProtoMessageUser,
		Movie:   movie,
		Privacy: messages.PRIVACY(theater.Privacy),
		VideoPlayerAccess: messages.PRIVACY(theater.VideoPlayerAccess),
	}, nil
}

func (s *Service) GetUserTheaters(ctx context.Context, req *proto.GetAllUserTheatersRequest) (*proto.UserTheatersResponse, error) {

	var (
		theaters   = make([]*messages.Theater, 0)
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

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	cursor, err := collection.Find(mCtx, bson.M{"user_id": authUser.ID})
	if err != nil {
		sentry.CaptureException(err)
		return &proto.UserTheatersResponse{
			Status:  "failed",
			Code:    http.StatusNotAcceptable,
			Message: "The requested parameter is not acceptable!",
		}, nil
	}

	for cursor.Next(mCtx) {
		theater := new(models.Theater)
		if err := cursor.Decode(theater); err != nil {
			sentry.CaptureException(err)
			break
		}
		th, err := SetDbTheaterToMessageTheater(mCtx, theater)
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