package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CastyLab/grpc.server/config"
	"golang.org/x/oauth2"
	"io/ioutil"
	"time"
)

const (
	userEndpoint  = "https://api.spotify.com/v1/me"
)

type JsonConfig struct {
	Web struct {
		ClientId        string   `json:"client_id"`
		ClientSecret    string   `json:"client_secret"`
		RedirectUri     string   `json:"redirect_uri"`
		AuthUri         string   `json:"auth_uri"`
		TokenUri        string   `json:"token_uri"`
	} `json:"web"`
}

var (
	err error
	jsonConfig []byte
	jsonConfigMap = new(JsonConfig)
	oauthClient *oauth2.Config
	scopes = []string{
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

func Configure() error {

	jsonConfig, err = ioutil.ReadFile(config.Map.Secrets.Oauth.Spotify)
	if err != nil {
		return fmt.Errorf("could not read spotify secret config file :%v", err)
	}

	if err := json.Unmarshal(jsonConfig, &jsonConfigMap); err != nil {
		return fmt.Errorf("could not parse spotify secret config file :%v", err)
	}

	oauthClient = &oauth2.Config{
		ClientID:     jsonConfigMap.Web.ClientId,
		ClientSecret: jsonConfigMap.Web.ClientSecret,
		Endpoint:     oauth2.Endpoint{
			AuthURL:  jsonConfigMap.Web.AuthUri,
			TokenURL: jsonConfigMap.Web.TokenUri,
		},
		RedirectURL:  jsonConfigMap.Web.RedirectUri,
		Scopes:       scopes,
	}

	return nil
}

func Authenticate(code string) (*oauth2.Token, error) {
	mCtx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	return oauthClient.Exchange(mCtx, code, oauth2.AccessTypeOnline)
}