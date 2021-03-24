package oauth

import (
	"fmt"

	"github.com/CastyLab/grpc.server/config"
	"github.com/CastyLab/grpc.server/oauth/google"
	"github.com/CastyLab/grpc.server/oauth/spotify"
)

func ConfigureOAUTHClients(c *config.ConfigMap) error {
	if err := google.Configure(c); err != nil {
		return fmt.Errorf("could not configure google oauth client")
	}
	if err := spotify.Configure(c); err != nil {
		return fmt.Errorf("could not configure spotify oauth client")
	}
	return nil
}
