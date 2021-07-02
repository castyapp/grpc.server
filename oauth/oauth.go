package oauth

import (
	"fmt"

	"github.com/castyapp/grpc.server/config"
	"github.com/castyapp/grpc.server/oauth/google"
	"github.com/castyapp/grpc.server/oauth/spotify"
)

func ConfigureOAUTHClients(c *config.Map) error {
	if err := google.Configure(c); err != nil {
		return fmt.Errorf("could not configure google oauth client")
	}
	if err := spotify.Configure(c); err != nil {
		return fmt.Errorf("could not configure spotify oauth client")
	}
	return nil
}
