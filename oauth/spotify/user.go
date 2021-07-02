package spotify

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

type ExplicitContent struct {
	FilterEnabled bool `json:"filter_enabled"`
	FilterLocked  bool `json:"filter_locked"`
}

type User struct {
	ID           string   `json:"id"`
	DisplayName  string   `json:"display_name"`
	Email        string   `json:"email"`
	Country      string   `json:"country"`
	Href         string   `json:"href"`
	Images       []string `json:"images"`
	Product      string   `json:"product"`
	Type         string   `json:"type"`
	URI          string   `json:"uri"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	}
}

func (u *User) GetUserID() string {
	return u.ID
}

func (u *User) GetAvatar() string {
	return ""
}

func (u *User) GetEmailAddress() string {
	return u.Email
}

func (u *User) GetFullname() string {
	return u.DisplayName
}

func GetUserByToken(token *oauth2.Token) (*User, error) {

	request, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var user *User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return user, nil
	}

	return nil, errors.New("could not get user from spotify")
}
