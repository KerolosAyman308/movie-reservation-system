package db

import (
	"fmt"
	log "log/slog"
	"time"

	"movie/system/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewConn attempts to connect to the DB with retry logic.
func NewConn(cfg config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	const maxRetries = 10

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: cfg.MysqlAddress,
		}), &gorm.Config{})

		if err == nil {
			return db, nil
		}

		log.Warn("Failed to connect to database", "attempt", i, "max", maxRetries, "error", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("database connection failed after %d attempts: %w", maxRetries, err)
}
