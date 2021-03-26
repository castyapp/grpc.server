package user

import (
	"context"
	"net/http"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetFriends(ctx context.Context, req *proto.AuthenticateRequest) (*proto.FriendsResponse, error) {

	user, err := auth.Authenticate(s.db, req)
	if err != nil {
		return nil, err
	}

	friends, err := helpers.GetFriendsFromDatabase(s.db, ctx, user)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not get friends!")
	}

	return &proto.FriendsResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: friends,
	}, nil
}
