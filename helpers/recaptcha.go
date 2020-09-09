package helpers

import (
	"encoding/json"
	"errors"
	"github.com/CastyLab/grpc.server/config"
	"net/http"
	"net/url"
	"strings"
)

type VerifyResponse struct {
	Success       bool       `json:"success"`
	ChallengeTs   string     `json:"challenge_ts"`
	Hostname      string     `json:"hostname"`
	ErrorCodes    []string   `json:"error-codes"`
}

func VerifyRecaptcha(code string) (bool, error) {

	verifyResp := new(VerifyResponse)

	params := url.Values{}
	params.Add("secret", config.Map.Secrets.HcaptchaSecret)
	params.Add("response", code)

	request, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", strings.NewReader(params.Encode()))
	if err != nil {
		return false, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&verifyResp); err != nil {
		return false, err
	}

	if verifyResp.Success == true {
		return true, nil
	}

	return false, errors.New("could not verify captcha")
}