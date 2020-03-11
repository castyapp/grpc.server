package auth

import (
	"context"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/CastyLab/grpc.server/oauth/google"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func (Service) CallbackOAUTH(ctx context.Context, req *proto.OAUTHRequest) (*proto.AuthResponse, error) {

	var (
		user         = new(models.User)
		collection   = db.Connection.Collection("users")
		mCtx, _      = context.WithTimeout(ctx, 10 * time.Second)
		unauthorized = &proto.AuthResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}
	)

	switch req.Service {
	case proto.OAUTHRequest_Google:

		token, err := google.Authenticate(req.Code)
		if err != nil {
			return unauthorized, err
		}

		oauthUser, err := google.GetUserByToken(token)
		if err != nil {
			return unauthorized, err
		}

		if err := collection.FindOne(mCtx, bson.M{ "email": oauthUser.Email }).Decode(&user); err != nil {
			sentry.CaptureException(err)
			return unauthorized, err
		}

	}

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