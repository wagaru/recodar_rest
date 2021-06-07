package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type LineResponseError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type LineTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	LineResponseError
}

type GoogleUserInfoResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type LineJWTClaims struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
	jwt.StandardClaims
}

func (delivery *httpDelivery) authLine(c *gin.Context) {
	state := utils.RandToken(30)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	log.Printf("Store state %s", state)

	request, err := http.NewRequest("GET", "https://access.line.me/oauth2/v2.1/authorize", nil)
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err})
	}

	query := request.URL.Query()
	query.Add("response_type", "code")
	query.Add("client_id", delivery.config.LineLoginClientID)
	query.Add("state", state)
	query.Add("redirect_uri", delivery.config.Server+delivery.config.LineLoginRedirectURL)
	query.Add("scope", "profile openid email")
	request.URL.RawQuery = query.Encode()
	c.Redirect(http.StatusFound, request.URL.String())
}

func (delivery *httpDelivery) authLineCallback(c *gin.Context) {
	errorStr, errorDescription := c.Query("error"), c.QueryArray("error_description")
	if errorStr != "" {
		delivery.WrapResponse(c, ErrorResponse{errMsg: errorStr, errDetail: strings.Join(errorDescription, ",")})
	}

	state, code := c.Query("state"), c.Query("code")

	// check state
	session := sessions.Default(c)
	log.Printf("Store state %s", session.Get("state"))
	if session.Get("state") != state {
		delivery.WrapResponse(c, ErrorResponse{errMsg: "Invalid State"})
	}

	resp, err := http.PostForm("https://api.line.me/oauth2/v2.1/token", url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {delivery.config.LineLoginClientID},
		"client_secret": {delivery.config.LineLoginClientSecret},
		"code":          {code},
		"redirect_uri":  {delivery.config.Server + delivery.config.LineLoginRedirectURL},
	})
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errMsg: "Invalid token"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		delivery.WrapResponse(c, ErrorResponse{errMsg: "Get Access token failed."})
	}

	respData := LineTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errDetail: "Decode resp data failed."})
	}

	if respData.Error != "" {
		delivery.WrapResponse(c, ErrorResponse{errMsg: respData.Error, errDetail: respData.ErrorDescription})
	}

	user := &domain.User{
		BindingSource: "line",
		AccessToken:   respData.AccessToken,
		RefreshToken:  respData.RefreshToken,
		// JWT:           respData.IDToken,
	}

	// parse JWT
	token, err := jwt.ParseWithClaims(respData.IDToken, &LineJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(delivery.config.LineLoginClientSecret), nil
	})

	if claims, ok := token.Claims.(*LineJWTClaims); ok && token.Valid {
		user.LineUserID = claims.StandardClaims.Subject
		user.Name = claims.Name
		user.Email = claims.Email
		user.Picture = claims.Picture
	} else {
		delivery.WrapResponse(c, ErrorResponse{err: err})
	}

	// update Data
	user, err = delivery.usecase.UpsertUser(context.Background(),
		map[string]interface{}{
			"binding_source": user.BindingSource,
			"line_user_id":   user.LineUserID,
		}, map[string]interface{}{
			"name":          user.Name,
			"email":         user.Email,
			"picture":       user.Picture,
			"access_token":  user.AccessToken,
			"refresh_token": user.RefreshToken,
			// "jwt":           user.JWT,
		})
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errMsg: "Upsert user failed."})
	}

	jwtToken, err := delivery.usecase.GenerateJWTToken(context.Background(), user)
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errMsg: "Generate token failed."})
	}
	delivery.WrapResponse(c, SuccessResponse{data: map[string]interface{}{
		"token": jwtToken,
	}})
}

func (delivery *httpDelivery) authGoogle(c *gin.Context) {
	// Ref: https://developers.google.com/identity/protocols/oauth2/openid-connect#php
	// One good choice for a state token is a string of 30 or so characters constructed using a high-quality random-number generator.
	state := utils.RandToken(30)
	conf := delivery.getGoogleOauthConfig()
	c.Redirect(http.StatusFound, conf.AuthCodeURL(state))
}

func (delivery *httpDelivery) authGoogleCallback(c *gin.Context) {
	errorStr := c.Query("error")
	if errorStr != "" {
		delivery.WrapResponse(c, ErrorResponse{errMsg: errorStr})
	}

	conf := delivery.getGoogleOauthConfig()
	// if session.Get("state") != c.Query("state") {
	// 	return "", errors.New("Invalid state")
	// }
	respToken, err := conf.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errDetail: "Google exchange token failed."})
	}
	token, refresh := respToken.AccessToken, respToken.RefreshToken
	log.Printf("token %s", token)

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errMsg: "Invalid"})
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errMsg: "Invalid"})
	}
	defer resp.Body.Close()

	userInfo := GoogleUserInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errMsg: "Invalid", errDetail: "Decode response failed."})
	}

	user := &domain.User{
		BindingSource: "google",
		AccessToken:   token,
		RefreshToken:  refresh,
		Email:         userInfo.Email,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
		GoogleUserID:  userInfo.ID,
	}

	user, err = delivery.usecase.UpsertUser(context.Background(),
		map[string]interface{}{
			"binding_source": user.BindingSource,
			"google_user_id": user.GoogleUserID,
		}, map[string]interface{}{
			"name":          user.Name,
			"email":         user.Email,
			"picture":       user.Picture,
			"access_token":  user.AccessToken,
			"refresh_token": user.RefreshToken,
		})
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errDetail: "Upsert user failed."})
	}

	log.Printf("user %v", user)

	jwtToken, err := delivery.usecase.GenerateJWTToken(context.Background(), user)
	if err != nil {
		delivery.WrapResponse(c, ErrorResponse{err: err, errDetail: "Generate JWT token failed"})
	}
	delivery.WrapResponse(c, SuccessResponse{data: map[string]interface{}{
		"token": jwtToken,
	}})
}

func (delivery *httpDelivery) getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     delivery.config.GoogleClientID,
		ClientSecret: delivery.config.GoogleClientSecret,
		RedirectURL:  delivery.config.Server + delivery.config.GoogleOauthRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
