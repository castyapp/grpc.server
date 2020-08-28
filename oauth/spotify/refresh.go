package spotify

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Token struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
}

func RefreshToken(refreshToken string) (*Token, error) {

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

	token := new(Token)
	if err := json.NewDecoder(response.Body).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}
