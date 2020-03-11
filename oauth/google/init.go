package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
)

var (
	err error
	jsonConfig []byte
	oauthClient *oauth2.Config
	scopes = []string{ "profile", "email", "openid" }
)

func init() {
	jsonConfig, err = ioutil.ReadFile("./oauth/google/client_secret.json")
	if err != nil {
		log.Fatal(err)
	}
	oauthClient, err = google.ConfigFromJSON(jsonConfig, scopes...)
	if err != nil {
		log.Fatal(err)
	}
}