package usecase

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/wagaru/recodar-rest/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (u *usecase) GetLineOAuthURL() string {
	state := utils.RandToken(30)
	request, err := http.NewRequest("GET", "https://access.line.me/oauth2/v2.1/authorize", nil)
	if err != nil {
		log.Printf("Generate Line Oauth URL failed: %v", err)
		return ""
	}

	query := request.URL.Query()
	query.Add("response_type", "code")
	query.Add("client_id", u.config.LineLoginClientID)
	query.Add("state", state)
	query.Add("redirect_uri", u.config.LineLoginRedirectURL)
	query.Add("scope", "profile openid email")
	request.URL.RawQuery = query.Encode()
	return request.URL.String()
}

func (u *usecase) GetGoogleOAuthURL() string {
	// Ref: https://developers.google.com/identity/protocols/oauth2/openid-connect#php
	// One good choice for a state token is a string of 30 or so characters constructed using a high-quality random-number generator.
	state := utils.RandToken(30)
	conf := u.getGoogleOauthConfig()
	return conf.AuthCodeURL(state)
}

func (u *usecase) GetGoogleOAuthAccessToken(state, code string) (string, string, time.Time, error) {
	conf := u.getGoogleOauthConfig()
	// if session.Get("state") != state {
	// 	return "", errors.New("Invalid state")
	// }
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return token.AccessToken, token.RefreshToken, token.Expiry, nil
}

func (u *usecase) getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     u.config.GoogleClientID,
		ClientSecret: u.config.GoogleClientSecret,
		RedirectURL:  u.config.GoogleOauthRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
