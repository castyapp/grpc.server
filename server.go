package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/db"
	"github.com/castyapp/grpc.server/jwt"
	"github.com/castyapp/grpc.server/oauth"
	"github.com/castyapp/grpc.server/redis"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/grpc.server/services/message"
	"github.com/castyapp/grpc.server/services/theater"
	"github.com/castyapp/grpc.server/services/user"
	"github.com/castyapp/grpc.server/storage"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configMap *config.ConfigMap
	server    *grpc.Server
	err       error
	port      *int
	host      *string
)

func init() {

	log.SetFlags(log.Ltime | log.Lshortfile)

	server = grpc.NewServer()
	configFileName := flag.String("config-file", "config.hcl", "config.hcl file")

	flag.Parse()
	log.Printf("Loading ConfigMap from file: [%s]", *configFileName)

	if configMap, err = config.LoadFile(*configFileName); err != nil {
		log.Fatal(fmt.Errorf("could not load config: %v", err))
	}

	if configMap.Sentry.Enabled {
		if err := sentry.Init(sentry.ClientOptions{Dsn: configMap.Sentry.Dsn}); err != nil {
			log.Fatal(fmt.Errorf("could not initilize sentry: %v", err))
		}
	}

	if err := redis.Configure(configMap); err != nil {
		log.Fatal(fmt.Errorf("could not configure redis : %v", err))
	}

	if err := jwt.Load(configMap); err != nil {
		err := fmt.Errorf("could not load jwt configuration: %v", err)
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	if err := oauth.ConfigureOAUTHClients(configMap); err != nil {
		err := fmt.Errorf("could not load oauth clients configurations: %v", err)
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	if err := storage.Configure(configMap); err != nil {
		err := fmt.Errorf("could not configure s3 bucket storage client: %v", err)
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	if err := db.Configure(configMap); err != nil {
		err := fmt.Errorf("could not configure mongodb client: %v", err)
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}

func main() {

	defer func() {

		// Since sentry emits events in the background we need to make sure
		// they are sent before we shut down
		sentry.Flush(time.Second * 5)

		if err := redis.Close(); err != nil {
			log.Println(fmt.Errorf("could not close redis connection: %v", err))
		}

	}()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(fmt.Errorf("could not create tcp listener: %v", err))
	}

	proto.RegisterAuthServiceServer(server, new(auth.Service))
	proto.RegisterUserServiceServer(server, new(user.Service))
	proto.RegisterTheaterServiceServer(server, new(theater.Service))
	proto.RegisterMessagesServiceServer(server, new(message.Service))

	reflection.Register(server)

	log.Println(fmt.Sprintf("Server running in tcp:%s:%d", configMap.Listener.Host, configMap.Listener.Port))
	if err := server.Serve(listener); err != nil {
		sentry.CaptureException(err)
		log.Fatal(fmt.Errorf("could not serve grpc.tcp.listener :%v", err))
	}

}
