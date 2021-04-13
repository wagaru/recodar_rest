package router

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/wagaru/Recodar/server/internal/domain"
)

type Router struct {
	*gin.Engine
	Middlewares map[string]gin.HandlerFunc
}

type Session interface {
	sessions.Session
}

func NewRouter() *Router {
	router := &Router{
		gin.Default(),
		newMiddlewares(),
	}
	router.Use(cors.Default())
	newSession(router)
	return router
}

func newSession(router *Router) {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	gob.Register(domain.SessionUser{})
}

func newMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"AuthRequired": func(c *gin.Context) {
			session := sessions.Default(c)
			user := session.Get("user")
			if user == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			c.Next()
		},
	}
}
