package oauth

import (
	"fmt"
	"github.com/CastyLab/grpc.server/oauth/discord"
	"github.com/CastyLab/grpc.server/oauth/google"
	"github.com/CastyLab/grpc.server/oauth/spotify"
)

func ConfigureOAUTHClients() error {
	if err := google.Configure(); err != nil {
		return fmt.Errorf("could not configure google oauth client")
	}
	if err := discord.Configure(); err != nil {
		return fmt.Errorf("could not configure discord oauth client")
	}
	if err := spotify.Configure(); err != nil {
		return fmt.Errorf("could not configure spotify oauth client")
	}
	return nil
}
