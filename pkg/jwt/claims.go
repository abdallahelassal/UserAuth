package jwt

import "github.com/golang-jwt/jwt/v4"


type JwtCustomClaims struct{
	UserID string `json:"user_id"`
	UserName string	`json:"user_name"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshToken struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

