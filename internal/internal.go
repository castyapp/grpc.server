package internal

import (
	"context"
	"github.com/CastyLab/grpc.server/internal/services/user"
	_ "github.com/joho/godotenv/autoload"
	"net"
	"net/http"
	"os"
)

type WebsocketInternalClient struct {
	http.Client
	UserService *user.InternalWsUserService
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
	)

	Client = &WebsocketInternalClient{
		Client: httpClient,
		UserService: userService,
	}
}