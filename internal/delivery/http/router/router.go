package router

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wagaru/recodar-rest/internal/config"
	"github.com/wagaru/recodar-rest/internal/domain"
)

type Router struct {
	*gin.Engine
	Middlewares map[string]gin.HandlerFunc
}

func NewRouter(config *config.Config) *Router {
	router := &Router{
		gin.Default(),
		newMiddlewares(config),
	}
	router.Use(cors.Default())
	return router
}

func newMiddlewares(config *config.Config) map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"AuthRequired": func(c *gin.Context) {
			tokens := strings.Split(c.GetHeader("Authorization"), " ")
			if len(tokens) != 2 || (len(tokens) == 2 && tokens[0] != "Bearer") {
				log.Printf("Incorrect bearer type.%v", tokens)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No permission"})
				return
			}
			token, err := jwt.ParseWithClaims(tokens[1], &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.JwtSecret), nil
			})
			if err != nil {
				log.Printf("decode failed.%v", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No permission"})
				return
			}

			if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
				c.Set("name", claims.Name)
				c.Set("email", claims.Email)
				// c.Set("picture", claims.Picture)
				c.Set("IDHex", claims.StandardClaims.Subject)
			}

			c.Next()
		},
	}
}
