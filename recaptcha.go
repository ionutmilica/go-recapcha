package recaptcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	verifyURL = "https://www.google.com/recaptcha/api/siteverify"
)

// Map some of the fields returned as json in the Recaptcha response
type ReCaptchaResponse struct {
	Success bool

	//missing-input-secret		The secret parameter is missing.
	//invalid-input-secret		The secret parameter is invalid or malformed.
	//missing-input-response	The response parameter is missing.
	//invalid-input-response	The response parameter is invalid or malformed.
	ErrorCodes []string `json:"error-codes"`
}

type ReCaptcha struct {
	// Secret key that you get from https://www.google.com/recaptcha
	SecretKey string

	// By default there is no timeout for the request
	Timeout time.Duration
}

func (r ReCaptcha) Verify(response string, remoteIp string) (*ReCaptchaResponse, error) {
	params := url.Values{}
	params.Set("secret", r.SecretKey)
	params.Set("remoteip", remoteIp)
	params.Set("response", response)

	client := &http.Client{Timeout: r.Timeout}
	resp, err := client.PostForm(verifyURL, params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &ReCaptchaResponse{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
