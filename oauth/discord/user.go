package discord

import (
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

type User struct {
	Id             string  `json:"id"`
	Username       string  `json:"username"`
	Verified       bool    `json:"verified"`
	Locale         string  `json:"locale"`
	MFAEnabled     bool    `json:"mfa_enabled"`
	Flags          int     `json:"flags"`
	Avatar         string  `json:"avatar"`
	Discriminator  string  `json:"discriminator"`
	Email          string  `json:"email"`
}

func (u *User) GetUserId() string {
	return u.Id
}

func (u *User) GetAvatar() string {
	return u.Avatar
}

func (u *User) GetEmailAddress() string {
	return u.Email
}

func (u *User) GetFullname() string {
	return u.Username
}

func GetUserByToken(token *oauth2.Token) (*User, error) {

	request, err := http.NewRequest("GET", "https://discordapp.com/api/v6/users/@me", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer " + token.AccessToken)

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
		user.Avatar = "https://cdn.discordapp.com/avatars/" + user.Id + "/" + user.Avatar + ".jpg"
		return user, nil
	}

	return nil, errors.New("could not get user from discord")
}