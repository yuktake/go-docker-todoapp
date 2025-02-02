package auth

import (
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwtv5.RegisteredClaims
}
