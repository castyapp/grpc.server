package google

import (
	"context"
	"time"

	"github.com/CastyLab/grpc.server/config"
	"golang.org/x/oauth2"
)

var (
	err         error
	jsonConfig  []byte
	oauthClient *oauth2.Config
	scopes      = []string{
		"profile",
		"email",
		"openid",
	}
)

func Configure(c *config.ConfigMap) error {
	oauthClient = &oauth2.Config{
		ClientID:     c.Oauth.Google.ClientID,
		ClientSecret: c.Oauth.Google.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.Oauth.Google.AuthUri,
			TokenURL: c.Oauth.Google.TokenUri,
		},
		RedirectURL: c.Oauth.Google.RedirectUri,
		Scopes:      scopes,
	}
	return nil
}

func Authenticate(code string) (*oauth2.Token, error) {
	mCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return oauthClient.Exchange(mCtx, code, oauth2.AccessTypeOffline)
}
