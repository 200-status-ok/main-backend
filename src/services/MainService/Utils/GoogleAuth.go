package Utils

import (
	"context"
	"errors"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
)

var googleClientId = utils.ReadFromEnvFile(".env", "GOOGLE_CLIENT_ID")
var googleClientSecret = utils.ReadFromEnvFile(".env", "GOOGLE_CLIENT_SECRET")

func GetGoogleOauthConfig() *oauth2.Config {
	redirectGoogleUrl := ""
	if os.Getenv("APP_ENV2") == "production" {
		redirectGoogleUrl = utils.ReadFromEnvFile(".env", "PRODUCTION_REDIRECT_GOOGLE_URL")
	} else {
		redirectGoogleUrl = utils.ReadFromEnvFile(".env", "LOCAL_REDIRECT_GOOGLE_URL")
	}
	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectGoogleUrl,
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return googleOauthConfig
}

func GetGoogleAuthURL(state string) string {
	return GetGoogleOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func GetGoogleUserInfo(code string, state string) ([]byte, error) {
	if state != "random-state" {
		return nil, errors.New("invalid oauth state")
	}
	token, err := GetGoogleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("code exchange wrong: " + err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.New("failed getting user info: " + err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("failed reading response body: " + err.Error())
	}
	return contents, nil
}
