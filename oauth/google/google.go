package google

import (
	"context"
	"time"

	"github.com/castyapp/grpc.server/config"
	"golang.org/x/oauth2"
)

var (
	oauthClient *oauth2.Config
	scopes      = []string{
		"profile",
		"email",
		"openid",
	}
)

func Configure(c *config.Map) error {
	oauthClient = &oauth2.Config{
		ClientID:     c.Oauth.Google.ClientID,
		ClientSecret: c.Oauth.Google.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.Oauth.Google.AuthURI,
			TokenURL: c.Oauth.Google.TokenURI,
		},
		RedirectURL: c.Oauth.Google.RedirectURI,
		Scopes:      scopes,
	}
	return nil
}

func Authenticate(code string) (*oauth2.Token, error) {
	mCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return oauthClient.Exchange(mCtx, code, oauth2.AccessTypeOffline)
}
