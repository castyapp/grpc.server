package auth

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/jwt"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	c  *config.ConfigMap
	db *mongo.Database
	proto.UnimplementedAuthServiceServer
}

func NewService(ctx context.Context) *Service {
	database := ctx.Value("db")
	if database == nil {
		log.Panicln("db value is required in context!")
	}
	configMap := ctx.Value("cm")
	if configMap == nil {
		log.Panicln("configMap value is required in context!")
	}
	return &Service{db: database.(*mongo.Database), c: configMap.(*config.ConfigMap)}
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
		collection   = s.db.Collection("users")
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

		token, refreshedToken, err := jwt.CreateNewTokens(s.db, ctx, user.ID.Hex())
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
