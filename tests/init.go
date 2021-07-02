package tests

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/jwt"
	"github.com/castyapp/grpc.server/providers"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/grpc.server/services/message"
	"github.com/castyapp/grpc.server/services/theater"
	"github.com/castyapp/grpc.server/services/user"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const configFileName = "./config_test.hcl"

func newContext() (*core.Context, error) {

	ctx := core.NewContext(context.Background())
	if err := ctx.Set("config.filepath", configFileName); err != nil {
		return nil, err
	}

	return ctx.With(

		// Registering configmap provider
		&providers.ConfigProvider{},

		// config database (mongodb)
		&providers.DatabaseProvider{},

		// configure jwt
		&providers.LambdaProvider{
			Registeration: func(ctx *core.Context) error {
				cm := ctx.MustGet("config.map").(*config.Map)
				if err := jwt.Load(cm); err != nil {
					return fmt.Errorf("could not load jwt configuration: %v", err)
				}
				return nil
			},
		},

		// configure redis connection
		&providers.RedisProvider{},
	), nil
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

func startGRPCServer() (*grpc.Server, *bufconn.Listener) {

	var (
		bufferSize = 1024 * 1024
		listener   = bufconn.Listen(bufferSize)
		server     = grpc.NewServer()
	)

	mockConext, err := newContext()
	if err != nil {
		log.Fatalf("could not create a new context: %v", err)
	}

	proto.RegisterAuthServiceServer(server, auth.NewService(mockConext))
	proto.RegisterUserServiceServer(server, user.NewService(mockConext))
	proto.RegisterTheaterServiceServer(server, theater.NewService(mockConext))
	proto.RegisterMessagesServiceServer(server, message.NewService(mockConext))

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()

	return server, listener
}
