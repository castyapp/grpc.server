package spotify

import (
	"context"
	"time"

	"github.com/castyapp/grpc.server/config"
	"golang.org/x/oauth2"
)

const (
	userEndpoint = "https://api.spotify.com/v1/me"
)

type JSONConfig struct {
	Web struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
		AuthURI      string `json:"auth_uri"`
		TokenURI     string `json:"token_uri"`
	} `json:"web"`
}

var (
	oauthClient *oauth2.Config
	scopes      = []string{
		"user-read-private",
		"user-read-email",
		"user-read-playback-state",
		"user-modify-playback-state",
		"user-library-read",
		"playlist-read-private",
		"streaming",
		"user-read-currently-playing",
	}
)

func Configure(c *config.Map) error {
	oauthClient = &oauth2.Config{
		ClientID:     c.Oauth.Spotify.ClientID,
		ClientSecret: c.Oauth.Spotify.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.Oauth.Spotify.AuthURI,
			TokenURL: c.Oauth.Spotify.TokenURI,
		},
		RedirectURL: c.Oauth.Spotify.RedirectURI,
		Scopes:      scopes,
	}
	return nil
}

func Authenticate(code string) (*oauth2.Token, error) {
	mCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return oauthClient.Exchange(mCtx, code, oauth2.AccessTypeOnline)
}
