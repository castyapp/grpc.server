package auth

import (
	"context"
	"net/http"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/castyapp/grpc.server/jwt"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.AuthResponse, error) {

	db, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	db = db.(*mongo.Database)

	if req.RefreshedToken == nil {
		return nil, errors.New("Refreshed token is required!")
	}

	newAuthToken, newRefreshedToken, err := jwt.RefreshToken(s.Context, string(req.RefreshedToken))
	if err != nil {
		sentry.CaptureException(err)
		return nil, errors.New("Could not create tokens, please try again later!")
	}

	return &proto.AuthResponse{
		Status:         "success",
		Code:           http.StatusOK,
		Token:          []byte(newAuthToken),
		RefreshedToken: []byte(newRefreshedToken),
	}, nil
}
