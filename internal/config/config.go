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
	Protocol           string
	HostName           string
	File               ConfigFile
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
		Protocol:           getValue("PROTOCOL", nil),
		HostName:           getValue("HOSTNAME", nil),
		File: ConfigFile{
			BucketName:    getValue("BUCKETNAME", nil),
			FilesBasePath: getValue("FILEBASEPATH", nil),
			AWSAccessKey:  getValue("AWS_ACCESS_KEY", nil),
			AWSSecretKey:  getValue("AWS_SECRET_KEY", nil),
			AWSHost:       getValue("AWS_URL", nil),
			UseFile:       getValueAsBool("USE_FILE", false),
		},
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

func getValueAsBool(key string, fallback bool) bool {
	str := strconv.FormatBool(fallback)
	value := getValue(key, &str)

	return value == "true"
}
