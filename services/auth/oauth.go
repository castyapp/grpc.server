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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"time"
)

func (Service) CallbackOAUTH(ctx context.Context, req *proto.OAUTHRequest) (*proto.AuthResponse, error) {

	var (
		user         = new(models.User)
		collection   = db.Connection.Collection("users")
		unauthorized = status.Error(codes.Unauthenticated, "Unauthorized!")
	)

	var oauthUser oauth.User
	switch req.Service {
	case proto.OAUTHRequest_Google:
		token, err := google.Authenticate(req.Code)
		if err != nil {
			return nil, unauthorized
		}
		oauthUser, err = google.GetUserByToken(token)
		if err != nil {
			return nil, unauthorized
		}
	case proto.OAUTHRequest_Discord:
		token, err := discord.Authenticate(req.Code)
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

	token, refreshedToken, err := jwt.CreateNewTokens(ctx, userObjectId)
	if err != nil {
		return nil, unauthorized
	}

	return &proto.AuthResponse{
		Status: "success",
		Code:   http.StatusOK,
		Token:  []byte(token),
		RefreshedToken:  []byte(refreshedToken),
	}, nil
}