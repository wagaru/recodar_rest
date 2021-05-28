package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wagaru/recodar-rest/internal/domain"
)

func (usecase *usecase) GenerateJWTToken(ctx context.Context, user *domain.User) (string, error) {
	now := time.Now()
	// jwtId := user.Email + strconv.FormatInt(time.Now().Unix(), 10)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		domain.Claims{
			Name:    user.Name,
			Email:   user.Email,
			Picture: user.Picture,
			StandardClaims: jwt.StandardClaims{
				// Audience:  user.Email,
				Subject:   user.ID.Hex(),
				ExpiresAt: now.Add(24 * time.Hour).Unix(),
				// Id:        jwtId,
				IssuedAt: now.Unix(),
				Issuer:   "Recodar",
				// NotBefore: now.Unix(),
			},
		})
	jwtSecret := []byte(usecase.config.JwtSecret)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("generate signed token failed: %w", err)
	}
	return token, nil
}
