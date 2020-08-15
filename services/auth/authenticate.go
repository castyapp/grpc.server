package auth

import (
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func Authenticate(req *proto.AuthenticateRequest) (user *models.User, err error) {

	if req == nil || req.Token == nil {
		return nil, status.Error(codes.InvalidArgument, "Token is required!")
	}

	stringToken := strings.ReplaceAll(string(req.Token), "Bearer ", "")

	user, err = jwt.DecodeAuthToken([]byte(stringToken))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized!")
	}

	return user, nil
}
