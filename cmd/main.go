package main

import (
	"fmt"
	log "log/slog"
	"net/http"
	"os"

	"movie/system/internal/config"
	"movie/system/internal/db"
)

func main() {
	cfg := config.Load()

	dbConn, err := db.NewConn(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Error("Failed to get underlying sql.DB", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	listenOn := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	log.Info("Server starting", "address", listenOn)
	if err := http.ListenAndServe(listenOn, InitializeAPI(dbConn, cfg)); err != nil {
		log.Error("Server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
