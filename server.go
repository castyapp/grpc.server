package main

import (
	"flag"
	"fmt"
	"github.com/CastyLab/grpc.proto/proto"
	_ "github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/CastyLab/grpc.server/services/message"
	"github.com/CastyLab/grpc.server/services/theater"
	"github.com/CastyLab/grpc.server/services/user"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	log.SetFlags(log.Ltime | log.Lshortfile)

	if err := sentry.Init(sentry.ClientOptions{ Dsn: os.Getenv("SENTRY_DSN") }); err != nil {
		log.Fatal(err)
	}

	// Since sentry emits events in the background we need to make sure
	// they are sent before we shut down
	defer sentry.Flush(time.Second * 5)

	var (
		server = grpc.NewServer()
		port   = flag.Int("port", 55283, "gRPC server port")
		host   = flag.String("host", "0.0.0.0", "gRPC server host")
	)

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	proto.RegisterAuthServiceServer(server, new(auth.Service))
	proto.RegisterUserServiceServer(server, new(user.Service))
	proto.RegisterTheaterServiceServer(server, new(theater.Service))
	proto.RegisterMessagesServiceServer(server, new(message.Service))

	reflection.Register(server)

	log.Println(fmt.Sprintf("Server running in tcp:%s:%d", *host, *port))
	if err := server.Serve(listener); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}