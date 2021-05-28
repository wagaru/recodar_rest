package domain

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Name    string
	Email   string
	Picture string
	jwt.StandardClaims
}
