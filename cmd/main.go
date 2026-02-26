package main

import (
	"fmt"
	"log"
	"movie/system/internal/db"
	"movie/system/internal/env"
	"net/http"
)

func main() {
	db, err := db.NewConn()
	if err != nil {
		log.Panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic(err)
	}

	defer sqlDB.Close()

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", env.Env.Port), InitializeAPI(db))
}
