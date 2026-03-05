package movies

import (
	"errors"
	"movie/system/pkg"
	"net/http"
)

var (
	ErrGenreDubName          = errors.New("genre with same name already exists")
	ErrGenreNotFound         = errors.New("this genre not found")
	ErrGenreHasMovies        = errors.New("this genre associated with other movies")
	ErrMovieDubName          = errors.New("movie with same name already exists")
	ErrMovieNotFound         = errors.New("movie not found")
	ErrSomeGenresNotFound    = errors.New("some genres associated with this movies is not found")
	ErrMovieOneGenre         = errors.New("movie must attach with at least one genre")
	ErrGenreDubNameAPI       = pkg.APIError[string]{StatusCode: http.StatusConflict, Code: "ErrGenreDubName", Message: "genre with same name already exists"}
	ErrMovieDubNameAPI       = pkg.APIError[string]{StatusCode: http.StatusConflict, Code: "ErrMovieDubName", Message: "movie with same name already exists"}
	ErrSomeGenresNotFoundAPI = pkg.APIError[string]{StatusCode: http.StatusNotFound, Code: "ErrSomeGenresNotFound", Message: "some genres associated with this movies is not found"}
	ErrMovieOneGenreAPI      = pkg.APIError[string]{StatusCode: http.StatusNotFound, Code: "ErrMovieOneGenre", Message: "movie must attach with at least one genre"}
)
