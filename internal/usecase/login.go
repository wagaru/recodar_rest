package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/wagaru/Recodar/server/internal/delivery/http/router"
	"github.com/wagaru/Recodar/server/internal/domain"
	"github.com/wagaru/Recodar/server/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (u *usecase) GetGoogleOauthURL(session router.Session) string {
	// Ref: https://developers.google.com/identity/protocols/oauth2/openid-connect#php
	// One good choice for a state token is a string of 30 or so characters constructed using a high-quality random-number generator.
	state := utils.RandToken(30)
	session.Set("state", state)
	session.Save()
	conf := u.getGoogleOauthConfig()
	return conf.AuthCodeURL(state)
}

func (u *usecase) GetGoogleOauthAccessToken(state, code string, session router.Session) (string, error) {
	conf := u.getGoogleOauthConfig()
	if session.Get("state") != state {
		return "", errors.New("Invalid state")
	}
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func (u *usecase) getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     u.config.GoogleClientID,
		ClientSecret: u.config.GoogleClientSecret,
		RedirectURL:  u.config.GoogleOauthRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
}

func (u *usecase) HandleUserLogin(session router.Session, info []byte, source string) error {
	userInfo := &domain.User{}
	err := json.Unmarshal(info, userInfo)
	if err != nil {
		log.Printf("Unmarshal error %v", err)
		return fmt.Errorf("Unmarshal info to User failed: %w", err)
	}

	if userInfo.Email == "" {
		return fmt.Errorf("no email data")
	}

	user, err := u.repo.GetUser(context.Background(), "email", userInfo.Email)
	if err != nil {
		return fmt.Errorf("Unmarshal info to User failed: %w", err)
	}
	if user.ID == primitive.NilObjectID {
		userInfo.CreatedAt = time.Now()
		userInfo.BindingSource = source
		u.repo.StoreUser(context.Background(), userInfo)
	} else {
		u.repo.UpdateUser(context.Background(), user.ID.Hex(), userInfo)
	}
	user, _ = u.repo.GetUser(context.Background(), "email", userInfo.Email)
	session.Set("user", domain.NewSessionUser(user.ID.Hex(), user.Name, user.Picture))
	session.Save()
	return nil
}
