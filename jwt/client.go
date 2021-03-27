package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	expireTimeInt,
	expireTimeRefreshedTokenInt int

	accessTokenSecret,
	refreshTokenSecret []byte
)

// read the key files before starting http handlers
func Load(c *config.ConfigMap) error {
	expireTimeInt = c.JWT.AccessToken.ExpiresAt.Value
	expireTimeRefreshedTokenInt = c.JWT.RefreshToken.ExpiresAt.Value
	accessTokenSecret = []byte(c.JWT.AccessToken.Secret)
	refreshTokenSecret = []byte(c.JWT.RefreshToken.Secret)
	return nil
}

func CreateNewTokens(ctx *core.Context, userid string) (token, refreshedToken string, err error) {
	//generate the auth token
	token, err = createAuthToken(userid)
	if err != nil {
		return
	}
	// generate the refresh token
	refreshedToken, err = createRefreshToken(ctx, userid)
	if err != nil {
		return
	}
	return
}

func createAuthToken(userid string) (token string, err error) {

	authTokenExp := time.Now().Add(time.Hour * time.Duration(expireTimeInt)).Unix()

	// create a signer for rsa 256
	authJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userid,
		ExpiresAt: authTokenExp,
	})

	// generate the auth token string
	token, err = authJwt.SignedString(accessTokenSecret)
	return
}

func createRefreshToken(ctx *core.Context, userid string) (refreshTokenString string, err error) {

	var userObjectId primitive.ObjectID
	userObjectId, err = primitive.ObjectIDFromHex(userid)
	if err != nil {
		return
	}

	refreshTokenExp := time.Now().Add(time.Hour * time.Duration(expireTimeRefreshedTokenInt))

	dbConn, err := ctx.Get("db.mongo")
	if err != nil {
		return "", err
	}

	var (
		db         = dbConn.(*mongo.Database)
		result     *mongo.InsertOneResult
		collection = db.Collection("refreshed_tokens")
	)

	result, err = collection.InsertOne(ctx, bson.M{
		"user_id":    userObjectId,
		"valid":      true,
		"created_at": time.Now(),
		"expires_at": refreshTokenExp,
	})
	if err != nil {
		return
	}

	resultID := result.InsertedID.(primitive.ObjectID)

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        resultID.Hex(),
		Subject:   userid,
		ExpiresAt: refreshTokenExp.Unix(),
	})

	// generate the refresh token string
	refreshTokenString, err = refreshJwt.SignedString(refreshTokenSecret)
	return
}

func checkRefreshToken(ctx *core.Context, id string) (*models.RefreshedToken, error) {

	dbConn, err := ctx.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db             = dbConn.(*mongo.Database)
		refreshedToken = new(models.RefreshedToken)
		collection     = db.Collection("refreshed_tokens")
	)

	refreshedTokenObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": refreshedTokenObjectId}).Decode(refreshedToken)
	if err != nil {
		return nil, err
	}

	if refreshedToken.ID.Hex() == id && refreshedToken.Valid {
		return refreshedToken, nil
	}

	return nil, errors.New("could not find refreshed token or maybe expired")
}

func RefreshToken(ctx *core.Context, refreshTokenString string) (token, refreshedToken string, err error) {

	var refreshToken *jwt.Token
	refreshToken, err = jwt.ParseWithClaims(refreshTokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return refreshTokenSecret, nil
	})

	if refreshToken == nil {
		err = errors.New("error reading jwt claims")
		return
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		err = errors.New("error reading jwt claims")
		return
	}

	if err = refreshTokenClaims.Valid(); err != nil {

		var refreshedTokenObjectId primitive.ObjectID
		refreshedTokenObjectId, err = primitive.ObjectIDFromHex(refreshTokenClaims.Id)
		if err != nil {
			return
		}

		if err = deleteRefreshToken(ctx, refreshedTokenObjectId); err != nil {
			return
		}

		return
	}

	dbRefreshedToken, rErr := checkRefreshToken(ctx, refreshTokenClaims.Id)
	if rErr != nil {
		err = errors.New("could not decode refresh token or maybe token expired")
		return
	}

	if refreshToken.Valid {
		if err = deleteRefreshToken(ctx, *dbRefreshedToken.ID); err != nil {
			return
		}
		token, refreshedToken, err = CreateNewTokens(ctx, dbRefreshedToken.UserId.Hex())
		return
	}

	err = errors.New("unauthorized")
	return
}

func deleteRefreshToken(ctx *core.Context, jti primitive.ObjectID) (err error) {

	dbConn, err := ctx.Get("db.mongo")
	if err != nil {
		return err
	}

	var (
		db         = dbConn.(*mongo.Database)
		result     *mongo.DeleteResult
		collection = db.Collection("refreshed_tokens")
	)

	result, err = collection.DeleteOne(ctx, bson.M{"_id": jti})
	if err != nil {
		return
	}

	if result.DeletedCount != 1 {
		err = errors.New("could not delete refresh token from db")
		return
	}

	return
}

func DecodeAuthToken(ctx *core.Context, token []byte) (user *models.User, err error) {

	database, err := ctx.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	db := database.(*mongo.Database)

	// now, check that it matches what's in the auth token claims
	var authToken *jwt.Token
	authToken, err = jwt.ParseWithClaims(string(token), &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return accessTokenSecret, nil
	})

	if authToken == nil || authToken.Claims == nil {
		err = errors.New("error reading jwt claims")
		return
	}

	authTokenClaims, ok := authToken.Claims.(*jwt.StandardClaims)
	if !ok || err != nil {
		err = errors.New("error reading jwt claims")
		return
	}

	if !authToken.Valid {
		err = errors.New("auth token is not valid")
		return
	}

	mCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(authTokenClaims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	usersCollection := db.Collection("users")
	user = new(models.User)

	if err := usersCollection.FindOne(mCtx, bson.M{"_id": objectId}).Decode(&user); err != nil {
		return nil, fmt.Errorf("invalid user")
	}

	return
}
