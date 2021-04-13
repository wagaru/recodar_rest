package http

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wagaru/Recodar/server/internal/delivery/http/router"
)

func (delivery *httpDelivery) googleLogin(c *gin.Context) {
	session := sessions.Default(c)
	url := delivery.usecase.GetGoogleOauthURL(session.(router.Session))
	c.Redirect(http.StatusMovedPermanently, url)
}

func (delivery *httpDelivery) googleLoginCallback(c *gin.Context) {
	session := sessions.Default(c)
	token, err := delivery.usecase.GetGoogleOauthAccessToken(c.Query("state"), c.Query("code"), session.(router.Session))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	userInfo, err := delivery.GetGoogleUserInfo(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	log.Printf("Fetch user info...%v", userInfo)
	err = delivery.usecase.HandleUserLogin(session.(router.Session), userInfo, "google")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.Redirect(http.StatusFound, "/")
}

func (delivery *httpDelivery) GetGoogleUserInfo(token string) ([]byte, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	log.Printf("userinfo: %v", string(bodyBytes))
	return bodyBytes, nil
}

func (delivery *httpDelivery) me(c *gin.Context) {
	session := sessions.Default(c)
	c.JSON(http.StatusOK, gin.H{"user": session.Get("user")})
}

func (delivery *httpDelivery) logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
