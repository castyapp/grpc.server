package spotify

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
)

func RefreshToken(refreshToken string) (*oauth2.Token, error) {

	params := url.Values{}
	params.Set("client_id", oauthClient.ClientID)
	params.Set("client_secret", oauthClient.ClientSecret)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	request, err := http.NewRequest(http.MethodPost, oauthClient.Endpoint.TokenURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	token := new(oauth2.Token)
	if err := json.NewDecoder(response.Body).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}
