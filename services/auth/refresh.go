package auth

import (
	"context"
	"errors"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/jwt"
	"net/http"
)

func (s *Service) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.AuthResponse, error) {

	if req.RefreshedToken == nil {
		return nil, errors.New("refreshed token is required")
	}

	newAuthToken, newRefreshedToken, err := jwt.RefreshToken(string(req.RefreshedToken))
	if err != nil {
		return &proto.AuthResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, err
	}

	return &proto.AuthResponse{
		Status: "success",
		Code:   http.StatusOK,
		Token:  []byte(newAuthToken),
		RefreshedToken:  []byte(newRefreshedToken),
	}, nil
}
