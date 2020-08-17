package helpers

import (
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/golang/protobuf/ptypes"
)

func NewProtoUser(user *models.User) (*proto.User, error) {

	lastLogin, _ := ptypes.TimestampProto(user.LastLogin)
	joinedAt,  _ := ptypes.TimestampProto(user.JoinedAt)
	updatedAt, _ := ptypes.TimestampProto(user.UpdatedAt)

	protoUser := &proto.User{
		Id:             user.ID.Hex(),
		Fullname:       user.Fullname,
		Username:       user.Username,
		Hash:           user.Hash,
		Email:          user.Email,
		IsActive:       user.IsActive,
		IsStaff:        user.IsStaff,
		Verified:       user.Verified,
		EmailVerified:  user.EmailVerified,
		Avatar:         user.Avatar,
		TwoFaEnabled:   user.TwoFactorAuthEnabled,
		TwoFaToken:     user.TwoFactorAuthToken,
		LastLogin:      lastLogin,
		JoinedAt:       joinedAt,
		UpdatedAt:      updatedAt,
	}

	return protoUser, nil
}