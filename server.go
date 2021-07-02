package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/jwt"
	"github.com/castyapp/grpc.server/oauth"
	"github.com/castyapp/grpc.server/providers"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/grpc.server/services/message"
	"github.com/castyapp/grpc.server/services/theater"
	"github.com/castyapp/grpc.server/services/user"
	"github.com/castyapp/grpc.server/storage"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ctx  *core.Context
	port *int
	host *string
)

func main() {

	log.SetFlags(log.Ltime | log.Lshortfile)

	configFileName := flag.String("config-file", "config.hcl", "config.hcl file")
	host = flag.String("host", "0.0.0.0", "grpc server host listener")
	port = flag.Int("port", 55283, "grpc server port listener")

	flag.Parse()
	log.Printf("Loading ConfigMap from file: [%s]", *configFileName)

	ctx = core.NewContext(context.Background())
	if err := ctx.Set("config.filepath", *configFileName); err != nil {
		log.Fatal(fmt.Errorf("could not set config.filepath to context: %v", err))
	}

	ctx.With(

		// Registering configmap provider
		&providers.ConfigProvider{},

		// Init sentry loggin if its enabled
		&providers.SentryProvider{},

		// config database (mongodb)
		&providers.DatabaseProvider{},

		// config redis connection
		&providers.RedisProvider{},

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

		// configure oauth clients
		&providers.LambdaProvider{
			Registeration: func(ctx *core.Context) error {
				cm := ctx.MustGet("config.map").(*config.Map)
				if err := oauth.ConfigureOAUTHClients(cm); err != nil {
					return fmt.Errorf("could not load oauth clients configurations: %v", err)
				}
				return nil
			},
		},

		// configure s3 bucket (minio) storage
		&providers.LambdaProvider{
			Registeration: func(ctx *core.Context) error {
				cm := ctx.MustGet("config.map").(*config.Map)
				if err := storage.Configure(cm); err != nil {
					return fmt.Errorf("could not configure s3 bucket storage client: %v", err)
				}
				return nil
			},
		},
	)

	defer ctx.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(fmt.Errorf("could not create tcp listener: %v", err))
	}

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, auth.NewService(ctx))
	proto.RegisterUserServiceServer(server, user.NewService(ctx))
	proto.RegisterTheaterServiceServer(server, theater.NewService(ctx))
	proto.RegisterMessagesServiceServer(server, message.NewService(ctx))

	reflection.Register(server)

	log.Println(fmt.Sprintf("Server running in tcp:%s:%d", *host, *port))
	if err := server.Serve(listener); err != nil {
		sentry.CaptureException(err)
		log.Fatal(fmt.Errorf("could not serve grpc.tcp.listener :%v", err))
	}

}
