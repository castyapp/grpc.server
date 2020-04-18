package jwt

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	expireTimeInt,
	expireTimeRefreshedTokenInt int

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey

	// location of the files used for signing and verification
	privKeyPath = os.Getenv("JWT_PRIVATE_KEY_PATH") // `$ openssl genrsa -out app.rsa 2048`
	pubKeyPath  = os.Getenv("JWT_PUBLIC_KEY_PATH") // `$ openssl rsa -in app.rsa -pubout > app.rsa.pub`

	// mongodb refreshed tokens collection
	usersCollection = db.Connection.Collection("users")
	collection = db.Connection.Collection("refreshed_tokens")
)

// read the key files before starting http handlers
func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}

	expireTimeRefreshedTokenInt, err = strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_VALID_TIME"))
	if err != nil {
		sentry.CaptureException(err)
	}

	expireTimeInt, err = strconv.Atoi(os.Getenv("JWT_EXPIRE_TIME"))
	if err != nil {
		sentry.CaptureException(err)
	}
}

func CreateNewTokens(ctx context.Context, userid string) (token, refreshedToken string, err error) {

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

	authClaims := jwt.StandardClaims{
		Subject: userid,
		ExpiresAt: authTokenExp,
	}

	// create a signer for rsa 256
	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), authClaims)

	// generate the auth token string
	token, err = authJwt.SignedString(signKey)
	return
}

func createRefreshToken(ctx context.Context, userid string) (refreshTokenString string, err error) {

	var userObjectId primitive.ObjectID
	userObjectId, err = primitive.ObjectIDFromHex(userid)
	if err != nil {
		return
	}

	refreshTokenExp := time.Now().Add(time.Hour * time.Duration(expireTimeRefreshedTokenInt))

	var result *mongo.InsertOneResult
	result, err = collection.InsertOne(ctx, bson.M{
		"user_id": userObjectId,
		"valid": true,
		"created_at": time.Now(),
		"expires_at": refreshTokenExp,
	})
	if err != nil {
		return
	}

	resultID := result.InsertedID.(primitive.ObjectID)
	refreshClaims := jwt.StandardClaims{
		Id: resultID.Hex(),
		Subject: userid,
		ExpiresAt: refreshTokenExp.Unix(),
	}

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// generate the refresh token string
	refreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func checkRefreshToken(ctx context.Context, id string) (*models.RefreshedToken, error) {

	refreshedToken := new(models.RefreshedToken)

	refreshedTokenObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{ "_id": refreshedTokenObjectId}).Decode(refreshedToken)
	if err != nil {
		return nil, err
	}

	if refreshedToken.ID.Hex() == id && refreshedToken.Valid {
		return refreshedToken, nil
	}
	
	return nil, errors.New("could not find refreshed token or maybe expired")
}

func RefreshToken(ctx context.Context, refreshTokenString string) (token, refreshedToken string, err error) {

	var refreshToken *jwt.Token
	refreshToken, err = jwt.ParseWithClaims(refreshTokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
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

func deleteRefreshToken(ctx context.Context, jti primitive.ObjectID) (err error) {

	var result *mongo.DeleteResult
	result, err = collection.DeleteOne(ctx, bson.M{ "_id": jti })
	if err != nil {
		return
	}

	if result.DeletedCount != 1 {
		err = errors.New("could not delete refresh token from db")
		return
	}

	return
}

func DecodeAuthToken(token []byte) (user *models.User, err error) {

	// now, check that it matches what's in the auth token claims
	var authToken *jwt.Token
	authToken, err = jwt.ParseWithClaims(string(token), &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
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

	ctx, _ := context.WithTimeout(context.Background(), 20 * time.Second)

	objectId, err := primitive.ObjectIDFromHex(authTokenClaims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	user = new(models.User)
	if err := usersCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user); err != nil {
		return nil, fmt.Errorf("invalid user")
	}

	return
}