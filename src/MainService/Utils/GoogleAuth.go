package Utils

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var googleClientId = ReadFromEnvFile(".env", "GOOGLE_CLIENT_ID")
var googleClientSecret = ReadFromEnvFile(".env", "GOOGLE_CLIENT_SECRET")

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "https://main-backend.iran.liara.run/api/v1/users/auth/google/callback",
	ClientID:     googleClientId,
	ClientSecret: googleClientSecret,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func GetGoogleAuthURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func GetGoogleUserInfo(code string, state string, c *gin.Context) ([]byte, error) {
	if state != "random-state" {
		return nil, errors.New("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(context.Background(), code)
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
