package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/CastyLab/grpc.server/services"
	"github.com/golang/protobuf/ptypes/any"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"
)

func (s *Service) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.AuthResponse, error)  {

	var (
		user     = req.User
		database = db.Connection
		existsUser = new(models.User)
		collection = database.Collection("users")
		validationErrors []*any.Any
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Captcha is required!")
	}

	recaptcha := md.Get("g-recaptcha-response")
	if success, err := helpers.VerifyRecaptcha(recaptcha[0]); err != nil || !success {
		return nil, status.Error(codes.InvalidArgument, "Captcha is required!")
	}

	_ = collection.FindOne(ctx, bson.M{ "username": user.Username }).Decode(existsUser)
	_ = collection.FindOne(ctx, bson.M{ "email": user.Email }).Decode(existsUser)

	if existsUser.Username == user.Username {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "username",
			Value: []byte("Username already exists!"),
		})
	}

	if existsUser.Email == user.Email {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "email",
			Value: []byte("Email already exists!"),
		})
	}

	if len(validationErrors) > 0 {
		return nil, status.ErrorProto(&spb.Status{
			Code: int32(codes.InvalidArgument),
			Message: "Validation Error!",
			Details: validationErrors,
		})
	}

	dbUser := bson.M{
		"fullname":   user.Fullname,
		"hash":       services.GenerateHash(),
		"username":   strings.ToLower(user.Username),
		"email":      user.Email,
		"password":   models.HashPassword(user.Password),
		"is_active":  true,
		"verified": false,
		"is_staff": false,
		"email_verified": false,
		"email_token": services.RandomString(40),
		"state":      int(proto.PERSONAL_STATE_OFFLINE),
		"activity":   bson.M{},
		"avatar":     "default",
		"last_login": time.Now(),
		"joined_at":  time.Now(),
		"updated_at": time.Now(),
	}

	result, err := collection.InsertOne(ctx, dbUser)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not create the user, Please try again later!")
	}

	resultID := result.InsertedID.(primitive.ObjectID)

	newAuthToken, newRefreshedToken, err := jwt.CreateNewTokens(ctx, resultID.Hex())
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not create the user, Please try again later!")
	}

	return &proto.AuthResponse{
		Status: "success",
		Code:   http.StatusOK,
		Token: []byte(newAuthToken),
		RefreshedToken: []byte(newRefreshedToken),
	}, nil
}
