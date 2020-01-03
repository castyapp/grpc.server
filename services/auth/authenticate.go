package auth

import (
	"errors"
	"log"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/jwt"
	"movie.night.gRPC.server/proto"
	"strings"
)

func Authenticate(req *proto.AuthenticateRequest) (user *models.User, err error) {

	if req.Token == nil {
		return nil, errors.New("token is required")
	}

	stringToken := strings.ReplaceAll(string(req.Token), "Bearer ", "")

	user, err = jwt.DecodeAuthToken([]byte(stringToken))
	if err != nil {
		log.Println(err)
		return nil, errors.New("unauthorized")
	}

	return user, nil
}
