package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *Service) GetFriends(ctx context.Context, req *proto.AuthenticateRequest) (*proto.FriendsResponse, error) {

	user, err := auth.Authenticate(req)
	if err != nil {
		return nil, err
	}

	friends, err := helpers.GetFriendsFromDatabase(ctx, user)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not get friends!")
	}

	return &proto.FriendsResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  friends,
	}, nil
}
