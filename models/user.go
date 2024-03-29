package models

import (
	"time"

	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID                   *primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Fullname             string               `bson:"fullname,omitempty" json:"fullname,omitempty"`
	Username             string               `bson:"username,omitempty" json:"username,omitempty"`
	Hash                 string               `bson:"user_hash,omitempty" json:"hash,omitempty"`
	Email                string               `bson:"email,omitempty" json:"email,omitempty"`
	Password             string               `bson:"password,omitempty" json:"-,omitempty"`
	Verified             bool                 `bson:"verified,omitempty" json:"verified,omitempty"`
	IsActive             bool                 `bson:"is_active,omitempty" json:"is_active,omitempty"`
	IsStaff              bool                 `bson:"is_staff,omitempty" json:"is_staff,omitempty"`
	EmailVerified        bool                 `bson:"email_verified,omitempty" json:"email_verified,omitempty"`
	EmailToken           string               `bson:"email_token,omitempty" json:"-"`
	TwoFactorAuthEnabled bool                 `bson:"two_fa_enabled,omitempty" json:"two_fa_enabled"`
	TwoFactorAuthToken   string               `bson:"two_fa_token,omitempty" json:"_"`
	State                proto.PERSONAL_STATE `bson:"state,omitempty" json:"state,omitempty"`
	Avatar               string               `bson:"avatar,omitempty" json:"avatar,omitempty"`
	RoleID               uint                 `bson:"role_id,omitempty" json:"role_id,omitempty"`
	LastLogin            time.Time            `bson:"last_login,omitempty" json:"last_login,omitempty"`
	JoinedAt             time.Time            `bson:"joined_at,omitempty" json:"joined_at,omitempty"`
	UpdatedAt            time.Time            `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (u *User) ToProto() *proto.User {
	lastLogin := timestamppb.New(u.LastLogin)
	joinedAt := timestamppb.New(u.JoinedAt)
	updatedAt := timestamppb.New(u.UpdatedAt)
	return &proto.User{
		Id:            u.ID.Hex(),
		Fullname:      u.Fullname,
		Username:      u.Username,
		Hash:          u.Hash,
		Email:         u.Email,
		IsActive:      u.IsActive,
		IsStaff:       u.IsStaff,
		Verified:      u.Verified,
		EmailVerified: u.EmailVerified,
		Avatar:        u.Avatar,
		TwoFaEnabled:  u.TwoFactorAuthEnabled,
		LastLogin:     lastLogin,
		JoinedAt:      joinedAt,
		UpdatedAt:     updatedAt,
	}
}

func (u *User) SetPassword(password string) {
	u.Password = HashPassword(password)
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
