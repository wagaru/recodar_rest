package http

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
	"github.com/wagaru/recodar-rest/internal/logger"
)

type Router struct {
	*gin.Engine
	Middlewares       map[string]gin.HandlerFunc
	connectionLimiter *ConnectionLimiter
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !router.connectionLimiter.GetConnection() {
		return
	}
	defer router.connectionLimiter.ReleaseConnection()
	router.Engine.ServeHTTP(w, r)
}

func NewRouter(config *config.Config) *Router {
	router := &Router{
		gin.Default(),
		newMiddlewares(config),
		newConnectionLimiter(100),
	}
	router.Use(cors.New(newCorsConfig()))
	router.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(config.SessionSecret))))
	return router
}

func newCorsConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8081", // swagger
	}
	config.AddAllowMethods("DELETE")
	config.AddAllowHeaders("Authorization")
	return config
}

var limiter = NewRateLimiters(10, 100)

func newMiddlewares(config *config.Config) map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"AuthRequired": func(c *gin.Context) {
			tokens := strings.Split(c.GetHeader("Authorization"), " ")
			if len(tokens) != 2 || (len(tokens) == 2 && tokens[0] != "Bearer") {
				WrapResponse(c, ErrorResponse{status: http.StatusUnauthorized, errMsg: "No permission", errDetail: "Incorrect bearer type."})
				return
			}
			token, err := jwt.ParseWithClaims(tokens[1], &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.JwtSecret), nil
			})
			if err != nil {
				WrapResponse(c, ErrorResponse{status: http.StatusUnauthorized, err: err, errMsg: "No permission"})
				return
			}

			if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
				c.Set("name", claims.Name)
				c.Set("email", claims.Email)
				// c.Set("picture", claims.Picture)
				c.Set("userID", claims.StandardClaims.Subject)
			}

			c.Next()
		},
		"RateLimit": func(c *gin.Context) {
			limiter := limiter.GetLimiter(c.Request.RemoteAddr)
			if !limiter.Allow() {
				c.AbortWithStatus(http.StatusTooManyRequests)
				logger.Logger.Print("Too many Requests")
				return
			}
			c.Next()
		},
	}
}
