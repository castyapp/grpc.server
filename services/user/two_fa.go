package user

import (
	"context"
	"net/http"
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/services"
	"github.com/castyapp/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GenerateRecoveryCodes(ctx context.Context, req *proto.AuthenticateRequest) (*proto.RecoveryCodesResponse, error) {

	var (
		codesColl      = s.db.Collection("users")
		failedResponse = status.Error(codes.Internal, "Could not generate recovery codes, Please try again later!")
	)

	user, err := auth.Authenticate(s.db, req)
	if err != nil {
		return nil, err
	}

	if user.TwoFactorAuthEnabled {
		return nil, status.Error(codes.Aborted, "Two-factor authentication already enabled for this user!")
	}

	recCodes := make([]interface{}, 0)
	protoRecCodes := make([]*proto.RecoveryCode, 0)

	for i := 0; i < 4; i++ {
		rc := &proto.RecoveryCode{Code: services.RandomString(4)}
		protoRecCodes = append(protoRecCodes, rc)
		recCodes = append(recCodes, bson.M{
			"code":       rc.Code,
			"user_id":    user.ID,
			"created_at": time.Now(),
		})
	}

	if _, err := codesColl.InsertMany(ctx, recCodes); err != nil {
		return nil, failedResponse
	}

	return &proto.RecoveryCodesResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: protoRecCodes,
	}, nil

}

func (s *Service) EnableTwoFactorAuth(ctx context.Context, req *proto.TwoFactorAuthRequest) (*proto.Response, error) {

	return nil, nil
}

func (s *Service) DisableTwoFactorAuth(ctx context.Context, req *proto.TwoFactorAuthRequest) (*proto.Response, error) {

	return nil, nil
}
