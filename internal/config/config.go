package config

import (
	log "log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               int
	MysqlAddress       string
	JWTSecret          string
	RefreshTokenSecret string
	JWTAudience        string
	JWTIssuer          string
	IsProduction       bool
}

// Load reads environment variables (with optional .env file) and returns a Config.
func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found, reading from environment variables")
	}

	defaultPort := 8000
	return Config{
		Port:               getValueAsInt("PORT", &defaultPort),
		MysqlAddress:       getValue("MYSQL_ADDRESS", nil),
		JWTSecret:          getValue("JWT_SECRET", nil),
		JWTAudience:        getValue("JWT_AUDIENCE", nil),
		JWTIssuer:          getValue("JWT_ISSUER", nil),
		RefreshTokenSecret: getValue("REFRESH_TOKEN_SECRET", nil),
		IsProduction:       getValue("APP_ENV", nil) == "production",
	}
}

func getValue(key string, fallback *string) string {
	value, exists := os.LookupEnv(key)
	if !exists && fallback != nil {
		return *fallback
	}
	return value
}

func getValueAsInt(key string, fallback *int) int {
	value, exists := os.LookupEnv(key)
	if !exists && fallback != nil {
		return *fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil && fallback != nil {
		return *fallback
	}
	return n
}
