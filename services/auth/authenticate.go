package auth

import (
	"errors"
	"gitlab.com/movienight1/grpc.proto"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/jwt"
	"strings"
)

func Authenticate(req *proto.AuthenticateRequest) (user *models.User, err error) {

	if req.Token == nil {
		return nil, errors.New("token is required")
	}

	stringToken := strings.ReplaceAll(string(req.Token), "Bearer ", "")

	user, err = jwt.DecodeAuthToken([]byte(stringToken))
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	return user, nil
}
