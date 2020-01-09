package auth

import (
	"context"
	"github.com/getsentry/sentry-go"
	"gitlab.com/movienight1/grpc.proto"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"movie.night.gRPC.server/db"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/jwt"
	"net/http"
	"regexp"
	"time"
)

type Service struct {}

func (s *Service) isEmail(user string) bool {

	re := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])" +
			"?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if re.MatchString(user) {
		return true
	}

	return false
}

func (s *Service) validatePassword(user *models.User, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	return err == nil
}

func (s *Service) Authenticate(ctx context.Context, req *proto.AuthRequest) (*proto.AuthResponse, error) {

	var (
		collection   = db.Connection.Collection("users")
		user         = new(models.User)
		mCtx, _      = context.WithTimeout(context.Background(), 20 * time.Second)
		unauthorized = &proto.AuthResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}
	)

	if req.User == "" {
		return unauthorized, nil
	}

	if req.Pass == "" {
		return &proto.AuthResponse{
			Status:  "failed",
			Code:    420,
			Message: "pass field is required",
		}, nil
	}

	var filter = bson.M{ "username": string(req.User) }
	if s.isEmail(req.User) {
		filter = bson.M{ "email": string(req.User) }
	}

	if err := collection.FindOne(mCtx, filter).Decode(&user); err != nil {
		sentry.CaptureException(err)
		return unauthorized, nil
	}

	if s.validatePassword(user, req.Pass) {

		token, refreshedToken, err := jwt.CreateNewTokens(user.ID.Hex())
		if err != nil {
			return unauthorized, err
		}

		return &proto.AuthResponse{
			Status: "success",
			Code:   http.StatusOK,
			Token:  []byte(token),
			RefreshedToken:  []byte(refreshedToken),
		}, nil
	}

	return unauthorized, nil
}

