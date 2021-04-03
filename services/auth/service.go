package auth

import (
	"context"
	"net/http"
	"regexp"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/jwt"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	*core.Context
	proto.UnimplementedAuthServiceServer
}

func NewService(ctx *core.Context) *Service {
	return &Service{Context: ctx}
}

func (s *Service) isEmail(user string) bool {
	re := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])" +
			"?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(user) {
		return true
	}
	return false
}

func ValidatePassword(user *models.User, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	return err == nil
}

func (s *Service) Authenticate(ctx context.Context, req *proto.AuthRequest) (*proto.AuthResponse, error) {

	var (
		db           = s.MustGet("db.mongo").(*mongo.Database)
		collection   = db.Collection("users")
		user         = new(models.User)
		unauthorized = status.Error(codes.Unauthenticated, "Unauthorized!")
	)

	if req.User == "" || req.Pass == "" {
		return nil, unauthorized
	}

	filter := bson.M{"username": req.User}
	if s.isEmail(req.User) {
		filter = bson.M{"email": req.User}
	}

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find user!")
	}

	if ValidatePassword(user, req.Pass) {

		token, refreshedToken, err := jwt.CreateNewTokens(s.Context, user.ID.Hex())
		if err != nil {
			sentry.CaptureException(err)
			return nil, status.Error(codes.Internal, "Could not create auth token, Please try again later!")
		}

		return &proto.AuthResponse{
			Status:         "success",
			Code:           http.StatusOK,
			Token:          []byte(token),
			RefreshedToken: []byte(refreshedToken),
		}, nil
	}

	return nil, unauthorized
}
