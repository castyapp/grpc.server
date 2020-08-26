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
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func (Service) CallbackOAUTH(ctx context.Context, req *proto.OAUTHRequest) (*proto.AuthResponse, error) {

	var (
		err            error
		authenticated  bool
		token          *oauth2.Token
		oauthUser      oauth.User
		user           = new(models.User)
		collection     = db.Connection.Collection("users")
		consCollection = db.Connection.Collection("connections")
		unauthorized   = status.Error(codes.Unauthenticated, "Unauthorized!")
	)

	if req.AuthRequest != nil {
		user, err = Authenticate(req.AuthRequest)
		if err != nil {
			return nil, err
		}
		authenticated = true
	}

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
		return nil, status.Error(codes.InvalidArgument, "Invalid oauth service")
	}

	if authenticated {
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
			return &proto.AuthResponse{
				Status:  "success",
				Code:    http.StatusOK,
				Message: "Connection created successfully!",
			}, nil
		}
	}

	var (
		connection = new(models.Connection)
		filter = bson.M{
			"service_user_id": oauthUser.GetUserId(),
		}
	)

	if err := consCollection.FindOne(ctx, filter).Decode(connection); err != nil {
		return nil, unauthorized
	}

	err = collection.FindOne(ctx, bson.M{ "_id": connection.UserId }).Decode(user)
	if err != nil {
		return nil, unauthorized
	}

	if user.ID != connection.UserId {
		return nil, unauthorized
	}

	authToken, refreshedToken, err := jwt.CreateNewTokens(ctx, user.ID.Hex())
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