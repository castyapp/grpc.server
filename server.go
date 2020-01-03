package main

import (
	"flag"
	"fmt"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	_ "movie.night.gRPC.server/db"
	"movie.night.gRPC.server/proto"
	"movie.night.gRPC.server/services/auth"
	"movie.night.gRPC.server/services/message"
	"movie.night.gRPC.server/services/theater"
	"movie.night.gRPC.server/services/user"
	"net"
	"os"
	"time"
)

func main() {

	log.SetFlags(log.Ltime | log.Lshortfile)

	//_ = os.Setenv("http_proxy", "socks5://127.0.0.1:65535")
	if err := sentry.Init(sentry.ClientOptions{ Dsn: os.Getenv("SENTRY_DNS") }); err != nil {
		log.Fatal(err)
	}

	// Since sentry emits events in the background we need to make sure
	// they are sent before we shut down
	defer sentry.Flush(time.Second * 5)

	var (
		server = grpc.NewServer()
		port   = flag.String("port", "55283", "gRPC server port")
		host   = flag.String("host", "0.0.0.0", "gRPC server host")
	)

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	proto.RegisterAuthServiceServer(server, new(auth.Service))
	proto.RegisterUserServiceServer(server, new(user.Service))
	proto.RegisterTheaterServiceServer(server, new(theater.Service))
	proto.RegisterMessagesServiceServer(server, new(message.Service))

	reflection.Register(server)

	log.Println(fmt.Sprintf("Server running in tcp:%s:%s", *host, *port))
	if err := server.Serve(listener); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

}