package google

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
)

type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (u *User) GetUserId() string {
	return u.ID
}

func (u *User) GetAvatar() string {
	return u.Picture
}

func (u *User) GetEmailAddress() string {
	return u.Email
}

func (u *User) GetFullname() string {
	return u.GivenName
}

func GetUserByToken(token *oauth2.Token) (user *User, err error) {
	user = new(User)
	httpClient := oauthClient.Client(context.Background(), token)
	response, err := httpClient.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(user); err != nil {
		return nil, err
	}
	return
}
