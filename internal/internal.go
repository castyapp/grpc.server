package internal

import (
	"context"
	"github.com/CastyLab/grpc.server/internal/services/theater"
	"github.com/CastyLab/grpc.server/internal/services/user"
	"net"
	"net/http"
	"os"
)

type WebsocketInternalClient struct {
	http.Client
	UserService    *user.InternalWsUserService
	TheaterService *theater.InternalWsTheaterService
}

var (
	Client *WebsocketInternalClient
)

func init() {

	var (
		address = os.Getenv("INTERNAL_UNIX_FILE")
		httpClient = http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", address)
				},
			},
		}
		// Internal websocket services
		userService = &user.InternalWsUserService{HttpClient: httpClient}
		theaterService = &theater.InternalWsTheaterService{HttpClient: httpClient}
	)

	Client = &WebsocketInternalClient{
		Client: httpClient,
		UserService: userService,
		TheaterService: theaterService,
	}
}