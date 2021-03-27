package tests

import (
	"context"
	"testing"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func mockUser() *proto.User {
	return &proto.User{
		Fullname: "test-user",
		Username: "go-test",
		Password: "random-password",
		Email:    "random-email@casty.test",
	}
}

func dropDatabase(t *testing.T) {
	db, err := mockConext.Get("db.mongo")
	assert.NoError(t, err)
	err = db.(*mongo.Database).Drop(context.TODO())
	assert.NoError(t, err)
}

func testGetUser(t *testing.T, name string, client proto.UserServiceClient, user *proto.User, token []byte) {
	t.Run(name, func(t *testing.T) {
		userResp, err := client.GetUser(context.TODO(), &proto.AuthenticateRequest{Token: token})
		assert.NoError(t, err)
		assert.Equal(t, userResp.Code, int64(200))
		assert.Equal(t, userResp.Status, "success")
		assert.Equal(t, userResp.Result.Username, user.Username)
		assert.Equal(t, userResp.Result.Email, user.Email)
		assert.Equal(t, userResp.Result.Fullname, user.Fullname)
	})
}

func TestAuthentication(t *testing.T) {

	dropDatabase(t)
	defer dropDatabase(t)

	ctx := context.TODO()
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(getBufDialer(grpcListener)), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	mockedUser := mockUser()

	t.Run("RegisterUser", func(t *testing.T) {
		client := proto.NewUserServiceClient(conn)
		resp, err := client.CreateUser(ctx, &proto.CreateUserRequest{User: mockedUser})
		assert.NoError(t, err)
		assert.Equal(t, resp.Code, int64(200))
		assert.Equal(t, resp.Status, "success")
		assert.NotEmpty(t, resp.Token)
		testGetUser(t, "GetRegisteredUser", client, mockedUser, resp.Token)
	})

	t.Run("LoginUser", func(t *testing.T) {

		userClient := proto.NewUserServiceClient(conn)
		authClient := proto.NewAuthServiceClient(conn)

		t.Run("WithUsernameAndPassword", func(t *testing.T) {
			resp, err := authClient.Authenticate(ctx, &proto.AuthRequest{
				User: mockedUser.Username,
				Pass: mockedUser.Password,
			})
			assert.NoError(t, err)
			assert.Equal(t, resp.Code, int64(200))
			assert.Equal(t, resp.Status, "success")
			assert.NotEmpty(t, resp.Token)
			testGetUser(t, "GetLoginUser", userClient, mockedUser, resp.Token)
		})

		t.Run("WithEmailAndPassword", func(t *testing.T) {
			resp, err := authClient.Authenticate(ctx, &proto.AuthRequest{
				User: mockedUser.Email,
				Pass: mockedUser.Password,
			})
			assert.NoError(t, err)
			assert.Equal(t, resp.Code, int64(200))
			assert.Equal(t, resp.Status, "success")
			assert.NotEmpty(t, resp.Token)
			testGetUser(t, "GetLoginUser", userClient, mockedUser, resp.Token)
		})

	})
}
