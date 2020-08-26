package auth

import (
	"context"
	"fmt"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/CastyLab/grpc.server/oauth"
	"github.com/CastyLab/grpc.server/oauth/discord"
	"github.com/CastyLab/grpc.server/oauth/google"
	"github.com/CastyLab/grpc.server/oauth/spotify"
	"github.com/CastyLab/grpc.server/services"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"time"
)

func (Service) CallbackOAUTH(ctx context.Context, req *proto.OAUTHRequest) (*proto.AuthResponse, error) {

	var (
		user           = new(models.User)
		collection     = db.Connection.Collection("users")
		consCollection = db.Connection.Collection("connections")
		unauthorized   = status.Error(codes.Unauthenticated, "Unauthorized!")
	)

	var (
		token     *oauth2.Token
		err       error
		oauthUser oauth.User
	)

	switch req.Service {
	case proto.Connection_SPOTIFY:
		token, err = spotify.Authenticate(req.Code)
		if err != nil {
			return nil, unauthorized
		}
		oauthUser, err = spotify.GetUserByToken(token)
		if err != nil {
			return nil, unauthorized
		}
	case proto.Connection_GOOGLE:
		token, err = google.Authenticate(req.Code)
		if err != nil {
			return nil, unauthorized
		}
		oauthUser, err = google.GetUserByToken(token)
		if err != nil {
			return nil, unauthorized
		}
	case proto.Connection_DISCORD:
		token, err = discord.Authenticate(req.Code)
		if err != nil {
			return nil, unauthorized
		}
		oauthUser, err = discord.GetUserByToken(token)
		if err != nil {
			return nil, unauthorized
		}
	default:
		return nil, unauthorized
	}

	filter := bson.M{
		"service_user_id": oauthUser.GetUserId(),
		"type":            req.Service,
		"user_id":         user.ID,
	}

	consCount, err := consCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "Could not create connection, Please try again later!")
	}

	if consCount == 0 {
		connection := bson.M{
			"service_user_id": oauthUser.GetUserId(),
			"name":            oauthUser.GetFullname(),
			"type":            req.Service,
			"access_token":    token.AccessToken,
			"refreshed_token": token.RefreshToken,
			"show_activity":   true,
			"user_id":         user.ID,
			"created_at":      time.Now(),
			"updated_at":      time.Now(),
		}
		if _, err := consCollection.InsertOne(ctx, connection); err != nil {
			sentry.CaptureException(fmt.Errorf("could not create connection :%v", err))
			return nil, status.Error(codes.Unavailable, "Could not create connection, Please try again later!")
		}
	}

	var (
		userObjectId string
		cursor = collection.FindOne(ctx, bson.M{
			"email": oauthUser.GetEmailAddress(),
		})
	)

	if err := cursor.Decode(&user); err != nil {

		avatarName, err := services.SaveAvatarFromUrl(oauthUser.GetAvatar())
		if err != nil {
			log.Println(err)
			avatarName = "default"
		}

		dbUser := bson.M{
			"fullname":   oauthUser.GetFullname(),
			"hash":       services.GenerateHash(),
			"username":   services.RandomUserName(),
			"email":      oauthUser.GetEmailAddress(),
			"password":   models.HashPassword(services.RandomString(20)),
			"is_active":  true,
			"verified": false,
			"is_staff": false,
			"email_verified": false,
			"email_token": services.RandomString(40),
			"state":      int(proto.PERSONAL_STATE_OFFLINE),
			"activity":   bson.M{},
			"avatar":     avatarName,
			"last_login": time.Now(),
			"joined_at":  time.Now(),
			"updated_at": time.Now(),
		}

		result, err := collection.InsertOne(ctx, dbUser)
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not create the user, Please try again later!")
		}
		userObjectId = result.InsertedID.(primitive.ObjectID).Hex()
	}

	if user.ID != nil {
		userObjectId = user.ID.Hex()
	}

	authToken, refreshedToken, err := jwt.CreateNewTokens(ctx, userObjectId)
	if err != nil {
		return nil, unauthorized
	}

	return &proto.AuthResponse{
		Status:          "success",
		Code:            http.StatusOK,
		Token:           []byte(authToken),
		RefreshedToken:  []byte(refreshedToken),
	}, nil
}