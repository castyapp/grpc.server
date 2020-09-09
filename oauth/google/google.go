package google

import (
	"context"
	"fmt"
	"github.com/CastyLab/grpc.server/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"time"
)

var (
	err error
	jsonConfig []byte
	oauthClient *oauth2.Config
	scopes = []string{
		"profile",
		"email",
		"openid",
	}
)

func Configure() error {
	jsonConfig, err = ioutil.ReadFile(config.Map.Secrets.Oauth.Google)
	if err != nil {
		return fmt.Errorf("could not read google secret config file :%v", err)
	}
	oauthClient, err = google.ConfigFromJSON(jsonConfig, scopes...)
	if err != nil {
		return fmt.Errorf("could not parse google secret config file :%v", err)
	}
	return nil
}

func Authenticate(code string) (*oauth2.Token, error) {
	mCtx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	return oauthClient.Exchange(mCtx, code, oauth2.AccessTypeOffline)
}