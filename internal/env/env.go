package env

import (
	log "log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type env struct {
	Port               int
	MysqlAddress       string
	JWTSecret          string
	RefreshTokenSecret string
	JWTAudience        string
	JWTIssuer          string
}

var Env env = initEnv()

func initEnv() env {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	defaultPort := 8000
	return env{
		Port:               getValueAsInt("PORT", &defaultPort),
		MysqlAddress:       getValue("MYSQLADDRESS", nil),
		JWTSecret:          getValue("JWTSecret", nil),
		JWTAudience:        getValue("JWTAudience", nil),
		JWTIssuer:          getValue("JWTIssuer", nil),
		RefreshTokenSecret: getValue("RefreshTokenSecret", nil),
	}
}

func getValue(key string, fallBack *string) string {
	value, isExist := os.LookupEnv(key)
	if !isExist && fallBack != nil {
		return *fallBack
	}

	return value
}

func getValueAsInt(key string, fallBack *int) int {
	value, isExist := os.LookupEnv(key)
	if !isExist && fallBack != nil {
		return *fallBack
	}

	valueAsInt, error := strconv.ParseInt(value, 10, 64)
	if error != nil && fallBack != nil {
		return *fallBack
	}

	return int(valueAsInt)
}
