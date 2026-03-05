package main

import (
	"fmt"
	"io"
	f "movie/system/internal/files"
	movies "movie/system/internal/movies/models"
	user "movie/system/internal/user"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

// NewConn attempts to connect to the DB and returns the DB instance or an error.
func main() {
	stmts, err := gormschema.New("mysql").Load(
		&user.User{},
		&movies.Genre{},
		&movies.Movie{},
		&movies.MovieGenres{},
		&f.File{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
