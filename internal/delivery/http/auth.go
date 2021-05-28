package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/domain"
)

func (delivery *httpDelivery) authLine(c *gin.Context) {
	c.Redirect(http.StatusFound, delivery.usecase.GetLineOAuthURL())
}

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

func (delivery *httpDelivery) authLineCallback(c *gin.Context) {
	errorStr, errorDescription := c.Query("error"), c.QueryArray("error_description")
	if errorStr != "" {
		log.Printf("line callback with error:%s, err_description: %s", errorStr, errorDescription)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User not grant permission."})
		return
	}

	state, code := c.Query("state"), c.Query("code")

	// check state
	log.Printf("state: %v", state)

	resp, err := http.PostForm("https://api.line.me/oauth2/v2.1/token", url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {delivery.config.LineLoginClientID},
		"client_secret": {delivery.config.LineLoginClientSecret},
		"code":          {code},
		"redirect_uri":  {delivery.config.LineLoginRedirectURL},
	})
	if err != nil {
		log.Printf("change token failed, %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unexpected error."})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("get access token failed.")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unexpected error."})
		return
	}

	respData := LineTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		log.Printf("decode resp data failed: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unexpected error."})
		return
	}

	if respData.Error != "" {
		log.Printf("error: %s, error_description: %s", respData.Error, respData.ErrorDescription)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unexpected error."})
		return
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
		log.Printf("parse JWT token failed, %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unexpected error."})
		return
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
		log.Printf("Upsert user failed. %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unexpected error."})
		return
	}

	jwtToken, err := delivery.usecase.GenerateJWTToken(context.Background(), user)
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

func (delivery *httpDelivery) authGoogle(c *gin.Context) {
	c.Redirect(http.StatusFound, delivery.usecase.GetGoogleOAuthURL())
}

func (delivery *httpDelivery) authGoogleCallback(c *gin.Context) {
	errorStr := c.Query("error")
	if errorStr != "" {
		log.Printf("google callback with error:%s", errorStr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User not grant permission."})
		return
	}
	token, refresh, _, err := delivery.usecase.GetGoogleOAuthAccessToken(c.Query("state"), c.Query("code"))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	log.Printf("token %s", token)

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid"})
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid"})
		return
	}
	defer resp.Body.Close()

	userInfo := GoogleUserInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
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
		log.Printf("Upsert user failed. %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unexpected error."})
		return
	}

	log.Printf("user %v", user)

	jwtToken, err := delivery.usecase.GenerateJWTToken(context.Background(), user)
	if err != nil {
		log.Printf("Generate JWT token failed.%v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}
