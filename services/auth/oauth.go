package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/castyapp/grpc.server/jwt"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/oauth"
	"github.com/castyapp/grpc.server/oauth/google"
	"github.com/castyapp/grpc.server/oauth/spotify"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CallbackOAUTH(ctx context.Context, req *proto.OAUTHRequest) (*proto.AuthResponse, error) {

	var (
		db             = s.MustGet("db.mongo").(*mongo.Database)
		err            error
		authenticated  bool
		token          *oauth2.Token
		oauthUser      oauth.User
		user           = new(models.User)
		collection     = db.Collection("users")
		consCollection = db.Collection("connections")
		unauthorized   = status.Error(codes.Unauthenticated, "Unauthorized!")
	)

	if req.AuthRequest != nil {
		if user, err = Authenticate(s.Context, req.AuthRequest); err != nil {
			return nil, unauthorized
		}
		authenticated = true
	}

	switch req.Service {
	case proto.Connection_SPOTIFY:
		token, err = spotify.Authenticate(req.Code)
		if err != nil {
			return nil, err
		}
		oauthUser, err = spotify.GetUserByToken(token)
		if err != nil {
			return nil, err
		}
	case proto.Connection_GOOGLE:
		token, err = google.Authenticate(req.Code)
		if err != nil {
			return nil, err
		}
		oauthUser, err = google.GetUserByToken(token)
		if err != nil {
			return nil, err
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "Invalid oauth service")
	}

	var (
		connection = new(models.Connection)
		filter     = bson.M{"service_user_id": oauthUser.GetUserID()}
	)

	if err = consCollection.FindOne(ctx, filter).Decode(connection); err != nil {
		if err == mongo.ErrNoDocuments {
			if authenticated {
				connection := bson.M{
					"service_user_id": oauthUser.GetUserID(),
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
		} else {
			return nil, err
		}
	}

	if authenticated {
		if connection.UserID != user.ID {
			return nil, status.Error(codes.AlreadyExists, "Connection already associated with another user!")
		}
		return nil, status.Error(codes.AlreadyExists, "Connection already exists!")
	}

	if err = collection.FindOne(ctx, bson.M{"_id": connection.UserID}).Decode(user); err != nil {
		return nil, err
	}

	authToken, refreshedToken, err := jwt.CreateNewTokens(s.Context, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	return &proto.AuthResponse{
		Status:         "success",
		Code:           http.StatusOK,
		Token:          []byte(authToken),
		RefreshedToken: []byte(refreshedToken),
	}, nil
}
