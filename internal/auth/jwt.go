package auth

import (
	"fmt"
	"time"

	"movie/system/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret        string
	refreshSecret string
	aud           string
	iss           string
}

// CustomClaims holds the JWT payload for access tokens.
type CustomClaims struct {
	UserID  uint `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewJWTAuthenticator(cfg config.Config) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:        cfg.JWTSecret,
		aud:           cfg.JWTAudience,
		iss:           cfg.JWTIssuer,
		refreshSecret: cfg.RefreshTokenSecret,
	}
}

func (a *JWTAuthenticator) GenerateTokenPair(userID uint, isAdmin bool) (string, string, error) {
	// 1. Create Access Token
	accessClaims := CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    a.iss,
			Audience:  jwt.ClaimStrings{a.aud},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(a.secret))
	if err != nil {
		return "", "", err
	}

	// 2. Create Refresh Token
	refreshClaims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		Issuer:    a.iss,
		Audience:  jwt.ClaimStrings{a.aud},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(a.refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateAccessToken parses and validates the access token
func (a *JWTAuthenticator) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, a.jwtCallBack,
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}

// ValidateRefreshToken parses and validates the refresh token
func (a *JWTAuthenticator) ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, a.jwtRefreshCallBack,
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}

func (a *JWTAuthenticator) jwtCallBack(t *jwt.Token) (any, error) {
	_, ok := t.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
	}
	return []byte(a.secret), nil
}

func (a *JWTAuthenticator) jwtRefreshCallBack(t *jwt.Token) (any, error) {
	_, ok := t.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
	}
	return []byte(a.refreshSecret), nil
}
