package auth

import (
	"context"
	"github.com/CastyLab/grpc.proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/CastyLab/grpc.server/oauth/discord"
	"github.com/CastyLab/grpc.server/oauth/google"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func (Service) CallbackOAUTH(ctx context.Context, req *proto.OAUTHRequest) (*proto.AuthResponse, error) {

	var (
		oauthEmail   string
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
			return unauthorized, nil
		}
		oauthUser, err := google.GetUserByToken(token)
		if err != nil {
			return unauthorized, nil
		}
		oauthEmail = oauthUser.Email

	case proto.OAUTHRequest_Discord:
		token, err := discord.Authenticate(req.Code)
		if err != nil {
			return unauthorized, nil
		}
		oauthUser, err := discord.GetUserByToken(token)
		if err != nil {
			return unauthorized, nil
		}
		oauthEmail = oauthUser.Email
	}

	cursor := collection.FindOne(mCtx, bson.M{ "email": oauthEmail })
	if err := cursor.Decode(&user); err != nil {
		return &proto.AuthResponse{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Could not find user!",
		}, nil
	}

	token, refreshedToken, err := jwt.CreateNewTokens(user.ID.Hex())
	if err != nil {
		return unauthorized, nil
	}

	return &proto.AuthResponse{
		Status: "success",
		Code:   http.StatusOK,
		Token:  []byte(token),
		RefreshedToken:  []byte(refreshedToken),
	}, nil
}