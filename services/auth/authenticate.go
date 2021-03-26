package auth

import (
	"strings"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Authenticate(db *mongo.Database, req *proto.AuthenticateRequest) (user *models.User, err error) {

	if req == nil || req.Token == nil {
		return nil, status.Error(codes.InvalidArgument, "Token is required!")
	}

	stringToken := strings.ReplaceAll(string(req.Token), "Bearer ", "")

	user, err = jwt.DecodeAuthToken(db, []byte(stringToken))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized!")
	}

	return user, nil
}
