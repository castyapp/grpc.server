package auth

import (
	"context"
	"net/http"

	"github.com/castyapp/grpc.server/jwt"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

func (s *Service) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.AuthResponse, error) {

	if req.RefreshedToken == nil {
		return nil, errors.New("refreshed token is required")
	}

	newAuthToken, newRefreshedToken, err := jwt.RefreshToken(s.Context, string(req.RefreshedToken))
	if err != nil {
		sentry.CaptureException(err)
		return nil, errors.New("could not create tokens, please try again later")
	}

	return &proto.AuthResponse{
		Status:         "success",
		Code:           http.StatusOK,
		Token:          []byte(newAuthToken),
		RefreshedToken: []byte(newRefreshedToken),
	}, nil
}
