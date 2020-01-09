package user

import (
	"context"
	"gitlab.com/movienight1/grpc.proto"
	"gitlab.com/movienight1/grpc.proto/messages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"movie.night.gRPC.server/db"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/jwt"
	"movie.night.gRPC.server/services"
	"net/http"
	"time"
)

func (s *Service) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.AuthResponse, error)  {

	var (
		user     = req.User
		database = db.Connection
		existsUser = new(models.User)
		collection = database.Collection("users")
		validationErrors []*proto.ValidationError
	)

	mCtx, _ := context.WithTimeout(ctx, 20 * time.Second)

	_ = collection.FindOne(mCtx, bson.M{ "username": user.Username }).Decode(existsUser)
	_ = collection.FindOne(mCtx, bson.M{ "email": user.Email }).Decode(existsUser)

	if existsUser.Username == user.Username {
		validationErrors = append(validationErrors, &proto.ValidationError{
			Field: "username",
			Errors: []string{
				"Username already exists!",
			},
		})
	}

	if existsUser.Email == user.Email {
		validationErrors = append(validationErrors, &proto.ValidationError{
			Field: "email",
			Errors: []string{
				"Email already exists!",
			},
		})
	}

	if len(validationErrors) > 0 {
		return &proto.AuthResponse{
			Status:  "failed",
			Message: "Validation Error!",
			ValidationError: validationErrors,
			Code:    420,
		}, nil
	}

	dbUser := bson.M{
		"fullname":   user.Fullname,
		"hash":       services.GenerateHash(),
		"username":   user.Username,
		"email":      user.Email,
		"password":   models.HashPassword(user.Password),
		"is_active":  true,
		"state":      int(messages.PERSONAL_STATE_OFFLINE),
		"activity":   bson.M{},
		"avatar":     "default",
		"last_login": time.Now(),
		"joined_at":  time.Now(),
		"updated_at": time.Now(),
	}

	result, err := collection.InsertOne(mCtx, dbUser)
	if err != nil {
		return &proto.AuthResponse{
			Status:  "failed",
			Message: "Could not create the user, Please try again later!",
			Code:    http.StatusInternalServerError,
		}, nil
	}

	resultID := result.InsertedID.(primitive.ObjectID)

	newAuthToken, newRefreshedToken, err := jwt.CreateNewTokens(resultID.Hex())
	if err != nil {
		return &proto.AuthResponse{
			Status:  "failed",
			Message: "Could not create auth token for user, Please try again later!",
			Code:    http.StatusInternalServerError,
		}, err
	}

	return &proto.AuthResponse{
		Status: "success",
		Code:   http.StatusOK,
		Token: []byte(newAuthToken),
		RefreshedToken: []byte(newRefreshedToken),
	}, nil
}
