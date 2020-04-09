package user

import (
	"encoding/json"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/pingcap/errors"
	"net/http"
	"net/url"
	"strings"
)

type InternalWsUserService struct {
	HttpClient http.Client
}

func (s *InternalWsUserService) SendNewNotificationsEvent(userId string) error {

	params := url.Values{}
	params.Set("user_id", userId)

	request, err := http.NewRequest("POST", "http://unix/user/@NewFriendRequestEvent", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := s.HttpClient.Do(request)
	if err != nil {
		return err
	}

	result := map[string] interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}

	if result["status"] == "success" {
		return nil
	}
	
	return errors.New("Something went wrong, Could not send event!")
}

func (s *InternalWsUserService) AcceptNotificationEvent(auth *proto.AuthenticateRequest, user *models.User, friendID string) error {

	protoUser, err := helpers.SetDBUserToProtoUser(user)
	if err != nil {
		return err
	}

	jsonUser, err := json.Marshal(protoUser)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("friend_id", friendID)
	params.Set("user", string(jsonUser))

	request, err := http.NewRequest("POST", "http://unix/user/@FriendRequestAcceptedEvent", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", string(auth.Token))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := s.HttpClient.Do(request)
	if err != nil {
		return err
	}

	result := map[string] interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}

	if result["status"] == "success" {
		return nil
	}

	return errors.New("Something went wrong, Could not send event!")
}
