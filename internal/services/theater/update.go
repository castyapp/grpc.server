package theater

import (
	"encoding/json"
	"github.com/pingcap/errors"
	"net/http"
	"net/url"
	"strings"
)

type InternalWsTheaterService struct {
	HttpClient http.Client
}

func (s *InternalWsTheaterService) SendMediaSourceUpdateEvent(theaterId, mediaSourceId string) error {
	params := url.Values{}
	params.Set("theater_id", theaterId)
	params.Set("media_source_id", mediaSourceId)

	request, err := http.NewRequest("POST", "http://unix/media_source/@updated", strings.NewReader(params.Encode()))
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

func (s *InternalWsTheaterService) SendTheaterUpdateEvent(theaterId string) error {
	params := url.Values{}
	params.Set("theater_id", theaterId)

	request, err := http.NewRequest("POST", "http://unix/theater/@updated", strings.NewReader(params.Encode()))
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
