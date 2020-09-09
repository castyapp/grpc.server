package main

import (
	"flag"
	"fmt"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/config"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/CastyLab/grpc.server/oauth"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/CastyLab/grpc.server/services/message"
	"github.com/CastyLab/grpc.server/services/theater"
	"github.com/CastyLab/grpc.server/services/user"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

var (
	server *grpc.Server
	port *int
	host *string
)

func init() {

	log.SetFlags(log.Ltime | log.Lshortfile)

	server = grpc.NewServer()
	port   = flag.Int("port", 55283, "gRPC server port")
	host   = flag.String("host", "0.0.0.0", "gRPC server host")
	configFileName := flag.String("config-file", "config.yml", "config.yaml file")

	flag.Parse()
	log.Printf("Loading ConfigMap from file: [%s]", *configFileName)

	if err := config.Load(*configFileName); err != nil {
		log.Fatal(fmt.Errorf("could not load config: %v", err))
	}

	if err := sentry.Init(sentry.ClientOptions{ Dsn: config.Map.Secrets.SentryDsn }); err != nil {
		log.Fatal(fmt.Errorf("could not initilize sentry: %v", err))
	}

	if err := jwt.Load(); err != nil {
		err := fmt.Errorf("could not load jwt configuration: %v", err)
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	if err := oauth.ConfigureOAUTHClients(); err != nil {
		err := fmt.Errorf("could not load oauth clients configurations: %v", err)
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	if err := db.Configure(); err != nil {
		err := fmt.Errorf("could not configure mongodb client: %v", err)
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}

func main() {

	// Since sentry emits events in the background we need to make sure
	// they are sent before we shut down
	defer sentry.Flush(time.Second * 5)

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

	log.Println(fmt.Sprintf("Server running in tcp:%s:%d", *host, *port))
	if err := server.Serve(listener); err != nil {
		sentry.CaptureException(err)
		log.Fatal(fmt.Errorf("could not serve grpc.tcp.listener :%v", err))
	}

}