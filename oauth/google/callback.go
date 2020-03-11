package google

import (
	"context"
	"golang.org/x/oauth2"
	"time"
)

func Authenticate(code string) (*oauth2.Token, error) {
	mCtx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	return oauthClient.Exchange(mCtx, code, oauth2.AccessTypeOffline)
}