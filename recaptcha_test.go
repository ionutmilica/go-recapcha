package recaptcha

import (
	"github.com/h2non/gock"
	"reflect"
	"testing"
)

func TestReCaptcha_Verify_Success(t *testing.T) {
	successResponse := map[string]interface{}{
		"success":      true,
		"challenge_ts": "2017-04-30T10:36:08Z",
		"hostname":     "localhost",
	}

	defer gock.Off()
	gock.New("https://www.google.com").
		Post("/recaptcha/api/siteverify").
		Reply(200).
		JSON(successResponse)

	captcha := ReCaptcha{
		SecretKey: "secret-key",
	}
	response, err := captcha.Verify("secret", "ip")
	if err != nil {
		t.Errorf("captcha.Verify shouldn't have returned error but got %s", err)
	}

	if response.Success != true {
		t.Error("Expected CaptchaResponse.Success to be true but got false")
	}

	if response.Hostname != "localhost" {
		t.Errorf("Expected the hostname to be `localhost`, but got `%s`", response.Hostname)
	}

	if response.ChallengeTs == nil {
		t.Error("Expected to receive ChallengeTs as Time object but got nil")
	}
}

func TestReCaptcha_Verify_InvalidInputResponse(t *testing.T) {
	errorResponse := map[string]interface{}{
		"success":     false,
		"error-codes": []string{"invalid-input-response"},
	}

	defer gock.Off()
	gock.New("https://www.google.com").
		Post("/recaptcha/api/siteverify").
		Reply(200).
		JSON(errorResponse)

	captcha := ReCaptcha{
		SecretKey: "secret-key",
	}
	response, err := captcha.Verify("secret", "ip")

	if err != nil {
		t.Errorf("captcha.Verify shouldn't have returned error but got %s", err)
	}

	if response.Success == true {
		t.Error("Expected CaptchaResponse.Success to be false but got true")
	}

	if !reflect.DeepEqual(response.ErrorCodes, []string{"invalid-input-response"}) {
		t.Errorf("Expected to get invalid-input-response error code, but got %s", response.ErrorCodes)
	}
}
