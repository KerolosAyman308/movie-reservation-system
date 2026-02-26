package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateTokenPair(userID uint, isAdmin bool) (string, string, error)
	ValidateAccessToken(tokenString string) (*jwt.Token, error)
	ValidateRefreshToken(tokenString string) (*jwt.Token, error)
}
