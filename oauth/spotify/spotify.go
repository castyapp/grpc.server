package spotify

import (
	"context"
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
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
	config = new(JsonConfig)
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

func init() {

	jsonConfig, err = ioutil.ReadFile("./oauth/spotify/client_secret.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(jsonConfig, &config); err != nil {
		log.Fatal(err)
	}

	oauthClient = &oauth2.Config{
		ClientID:     config.Web.ClientId,
		ClientSecret: config.Web.ClientSecret,
		Endpoint:     oauth2.Endpoint{
			AuthURL:  config.Web.AuthUri,
			TokenURL: config.Web.TokenUri,
		},
		RedirectURL:  config.Web.RedirectUri,
		Scopes:       scopes,
	}

}

func Authenticate(code string) (*oauth2.Token, error) {
	mCtx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	return oauthClient.Exchange(mCtx, code, oauth2.AccessTypeOnline)
}