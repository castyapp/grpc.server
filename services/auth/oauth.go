package auth

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/CastyLab/grpc.server/oauth"
	"github.com/CastyLab/grpc.server/oauth/discord"
	"github.com/CastyLab/grpc.server/oauth/google"
	"github.com/CastyLab/grpc.server/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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

	var oauthUser oauth.User
	switch req.Service {
	case proto.OAUTHRequest_Google:
		token, err := google.Authenticate(req.Code)
		if err != nil {
			return unauthorized, nil
		}
		oauthUser, err = google.GetUserByToken(token)
		if err != nil {
			return unauthorized, nil
		}
	case proto.OAUTHRequest_Discord:
		token, err := discord.Authenticate(req.Code)
		if err != nil {
			return unauthorized, nil
		}
		oauthUser, err = discord.GetUserByToken(token)
		if err != nil {
			return unauthorized, nil
		}
	default:
		return unauthorized, nil
	}

	var (
		userObjectId string
		cursor = collection.FindOne(mCtx, bson.M{ "email": oauthUser.GetEmailAddress() })
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

		result, err := collection.InsertOne(mCtx, dbUser)
		if err != nil {
			return &proto.AuthResponse{
				Status:  "failed",
				Message: "Could not create the user, Please try again later!",
				Code:    http.StatusInternalServerError,
			}, nil
		}
		userObjectId = result.InsertedID.(primitive.ObjectID).Hex()
	}

	if user.ID != nil {
		userObjectId = user.ID.Hex()
	}

	token, refreshedToken, err := jwt.CreateNewTokens(mCtx, userObjectId)
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