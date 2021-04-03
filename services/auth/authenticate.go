package auth

import (
	"strings"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Authenticate(ctx *core.Context, req *proto.AuthenticateRequest) (user *models.User, err error) {

	if req == nil || req.Token == nil {
		return nil, status.Error(codes.InvalidArgument, "Token is required!")
	}

	stringToken := strings.ReplaceAll(string(req.Token), "Bearer ", "")

	user, err = jwt.DecodeAuthToken(ctx, []byte(stringToken))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized!")
	}

	return user, nil
}
