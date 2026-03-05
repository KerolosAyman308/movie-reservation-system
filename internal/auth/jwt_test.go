package auth

import (
	"movie/system/internal/config"
	"testing"
)

func TestJWTAuthenticator(t *testing.T) {
	cfg := config.Config{
		JWTSecret:          "secret",
		RefreshTokenSecret: "refresh-secret",
		JWTAudience:        "aud",
		JWTIssuer:          "iss",
	}

	auth := NewJWTAuthenticator(cfg)

	t.Run("Generate and validate token pair", func(t *testing.T) {
		userID := uint(1)
		isAdmin := true

		accessToken, refreshToken, err := auth.GenerateTokenPair(userID, isAdmin)
		if err != nil {
			t.Fatalf("Failed to generate token pair: %v", err)
		}
		if accessToken == "" || refreshToken == "" {
			t.Errorf("Expected token string, got empty")
		}

		// Validate Access Token
		accToken, err := auth.ValidateAccessToken(accessToken)
		if err != nil {
			t.Fatalf("Failed to validate access token: %v", err)
		}
		if !accToken.Valid {
			t.Errorf("Access token is not valid")
		}

		// Validate Refresh Token
		refToken, err := auth.ValidateRefreshToken(refreshToken)
		if err != nil {
			t.Fatalf("Failed to validate refresh token: %v", err)
		}
		if !refToken.Valid {
			t.Errorf("Refresh token is not valid")
		}
	})
}
