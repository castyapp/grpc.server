package auth

import (
	"context"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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
		mCtx, _      = context.WithTimeout(ctx, 20 * time.Second)
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

	var filter = bson.M{ "username": req.User }
	if s.isEmail(req.User) {
		filter = bson.M{ "email": req.User }
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

