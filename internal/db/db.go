package db

import (
	"fmt"
	"log"
	"movie/system/internal/env"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewConn attempts to connect to the DB and returns the DB instance or an error.
func NewConn() (*gorm.DB, error) {
	dsn := env.Env.MysqlAddress

	var db *gorm.DB
	var err error
	var maxRetries int = 10
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: dsn,
		}), &gorm.Config{})

		//db.AutoMigrate(&user.User{})
		if err == nil {
			// Return the connection
			return db, nil
		}

		log.Printf("Attempt %d/%d failed to connect to database: %v", i, maxRetries, err)

		// Wait 2 seconds before trying again
		time.Sleep(2 * time.Second)
	}

	// Return the final error
	return nil, fmt.Errorf("database connection failed after %d attempts: %w", maxRetries, err)
}
