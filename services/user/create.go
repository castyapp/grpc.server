package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/CastyLab/grpc.server/services"
	"github.com/getsentry/sentry-go"
	"github.com/golang/protobuf/ptypes/any"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var invalidUsernames = []string{
	"login",
	"logout",
	"register",
	"iforgot",
	"settings",
	"messages",
	"home",
	"me",
	"profile",
	"callback",
	"oauth",
	"terms",
}

func (s *Service) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.AuthResponse, error) {

	var (
		user             = req.User
		database         = db.Connection
		validationErrors []*any.Any
		existsUser       = new(models.User)
		collection       = database.Collection("users")
		thCollection     = database.Collection("theaters")
	)

	for _, invalid := range invalidUsernames {
		if req.User.Username == invalid {
			return nil, status.ErrorProto(&spb.Status{
				Code:    int32(codes.InvalidArgument),
				Message: "Validation Error!",
				Details: []*any.Any{
					{
						TypeUrl: "username",
						Value:   []byte("Username is not available!"),
					},
				},
			})
		}
	}

	if strings.Contains(user.Username, "/") {
		return nil, status.ErrorProto(&spb.Status{
			Code:    int32(codes.InvalidArgument),
			Message: "Validation Error!",
			Details: []*any.Any{
				{
					TypeUrl: "username",
					Value:   []byte("Username is not available!"),
				},
			},
		})
	}

	if err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(existsUser); err != nil {
		sentry.CaptureException(err)
	}

	if err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(existsUser); err != nil {
		sentry.CaptureException(err)
	}

	if existsUser.Username == user.Username {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "username",
			Value:   []byte("Username already exists!"),
		})
	}

	if existsUser.Email == user.Email {
		validationErrors = append(validationErrors, &any.Any{
			TypeUrl: "email",
			Value:   []byte("Email already exists!"),
		})
	}

	if len(validationErrors) > 0 {
		return nil, status.ErrorProto(&spb.Status{
			Code:    int32(codes.InvalidArgument),
			Message: "Validation Error!",
			Details: validationErrors,
		})
	}

	dbUser := bson.M{
		"fullname":       user.Fullname,
		"hash":           services.GenerateHash(),
		"username":       strings.ToLower(user.Username),
		"email":          user.Email,
		"password":       models.HashPassword(user.Password),
		"is_active":      true,
		"verified":       false,
		"is_staff":       false,
		"email_verified": false,
		"email_token":    services.RandomString(40),
		"state":          int(proto.PERSONAL_STATE_OFFLINE),
		"two_fa_enabled": false,
		"two_fa_token":   fmt.Sprintf("re_token_%s", services.RandomString(30)),
		"avatar":         "default",
		"last_login":     time.Now(),
		"joined_at":      time.Now(),
		"updated_at":     time.Now(),
	}

	result, err := collection.InsertOne(ctx, dbUser)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Could not create the user, Please try again later!")
	}

	resultID := result.InsertedID.(primitive.ObjectID)

	newAuthToken, newRefreshedToken, err := jwt.CreateNewTokens(ctx, resultID.Hex())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Could not create the user, Please try again later!")
	}

	theater := bson.M{
		"description":         fmt.Sprintf("%s's Theater", dbUser["fullname"]),
		"privacy":             proto.PRIVACY_PUBLIC,
		"video_player_access": proto.VIDEO_PLAYER_ACCESS_ACCESS_BY_USER,
		"user_id":             resultID,
		"created_at":          time.Now(),
		"updated_at":          time.Now(),
	}

	_, err = thCollection.InsertOne(ctx, theater)
	if err != nil {
		sentry.CaptureException(fmt.Errorf("could not create user!: %v", err))
		_, err := collection.DeleteOne(ctx, bson.M{"_id": resultID})
		if err != nil {
			sentry.CaptureException(fmt.Errorf("could not failed user's deletation!: %v", err))
		}
		return nil, status.Error(codes.Internal, "Could not create user! Please try again later!")
	}

	return &proto.AuthResponse{
		Status:         "success",
		Code:           http.StatusOK,
		Token:          []byte(newAuthToken),
		RefreshedToken: []byte(newRefreshedToken),
	}, nil
}
