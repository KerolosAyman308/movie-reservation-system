package main

import (
	"fmt"
	log "log/slog"
	"movie/system/internal/db"
	"movie/system/internal/env"
	"net/http"
	"os"
)

func main() {
	db, err := db.NewConn()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	defer sqlDB.Close()

	listenOn := fmt.Sprintf("0.0.0.0:%d", env.Env.Port)
	http.ListenAndServe(listenOn, InitializeAPI(db))
	log.Info("Server started listening on $w", listenOn)
}
