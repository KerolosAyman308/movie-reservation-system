package movies

import (
	"errors"
	log "log/slog"
	e "movie/system/internal/movies"
	d "movie/system/internal/movies/dtos"
	s "movie/system/internal/movies/services"
	"movie/system/pkg"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MoviesHandler struct {
	GenreService  s.IGenreService
	MovieService  s.IMoviesService
	MaxFileSizeMB int64
}

func NewMovieHandler(genreService s.IGenreService, movieService s.IMoviesService) *MoviesHandler {
	return &MoviesHandler{
		GenreService: genreService,
		MovieService: movieService,
	}
}

func (h *MoviesHandler) FindPaginated(w http.ResponseWriter, r *http.Request) {
	query := pkg.NewPaginationRequest(r)
	genres, err := h.GenreService.GetPaginated(r.Context(), query)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	pkg.Ok(genres, "Genres retrieved successfully", w)
}

func (h *MoviesHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var dto d.CreateMovieDTO
	err := pkg.ReadJSON(w, r, &dto)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	createdGenre, err := h.MovieService.Create(r.Context(), &dto)
	if err != nil {
		if errors.Is(e.ErrMovieDubName, err) {
			pkg.Error(e.ErrMovieDubNameAPI, w, r)
			return
		}
		if errors.Is(e.ErrSomeGenresNotFound, err) {
			pkg.Error(e.ErrSomeGenresNotFoundAPI, w, r)
			return
		}
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok(createdGenre, "Created Movie successfully", w)
}

func (h *MoviesHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(h.MaxFileSizeMB << 20)
	movieId, err := strconv.ParseUint(chi.URLParam(r, "movieId"), 10, 64)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	defer file.Close()

	imageURL, err := h.MovieService.UploadImage(r.Context(), movieId, file, handler.Filename)
	if err != nil {
		if errors.Is(err, e.ErrMovieNotFound) {
			pkg.NotFound(w, r, (*any)(nil))
			return
		}
		log.Error(err.Error(), "Location", "Image upload movies")
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok(imageURL, "Created Movie successfully", w)
}

func (h *MoviesHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	var dto d.CreateGenreDTO
	err := pkg.ReadJSON(w, r, &dto)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	createdGenre, err := h.GenreService.Create(r.Context(), &dto)
	if err != nil {
		if errors.Is(e.ErrGenreDubName, err) {
			pkg.Error(e.ErrGenreDubNameAPI, w, r)
			return
		}
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok(createdGenre, "Created genre successfully", w)
}

func (h *MoviesHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	genreId, err := strconv.ParseUint(chi.URLParam(r, "genreId"), 10, 64)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	err = h.GenreService.Delete(r.Context(), genreId)
	if err != nil {
		if errors.Is(err, e.ErrGenreNotFound) {
			pkg.NotFound(w, r, (*any)(nil))
			return
		}
		if errors.Is(err, e.ErrGenreHasMovies) {
			pkg.BadRequestWithCustomMessage(w, r, e.ErrGenreHasMovies.Error(), (*any)(nil))
			return
		}
	}

	pkg.Ok((*any)(nil), "Delete successfully", w)
}

func (h *MoviesHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.ParseUint(chi.URLParam(r, "movieId"), 10, 64)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	err = h.MovieService.Delete(r.Context(), movieId)
	if err != nil {
		if errors.Is(err, e.ErrMovieNotFound) {
			pkg.NotFound(w, r, (*any)(nil))
			return
		}
		pkg.InternalError(w, r, (*any)(nil))
		return
	}

	pkg.Ok((*any)(nil), "Delete successfully", w)
}

func (h *MoviesHandler) MoviePaginated(w http.ResponseWriter, r *http.Request) {
	query := pkg.NewPaginationRequest(r)

	// Extract search parameters from request
	searchByTitle := r.URL.Query().Get("title")
	searchByName := r.URL.Query().Get("name")
	searchByGenre := r.URL.Query().Get("genre")

	movies, err := h.MovieService.GetPaginated(
		r.Context(),
		query,
		searchByTitle,
		searchByName,
		searchByGenre,
	)

	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	pkg.Ok(movies, "Movies retrieved successfully", w)
}

// AddGenres adds multiple genres to a movie
func (h *MoviesHandler) AddMovieGenres(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.ParseUint(chi.URLParam(r, "movieId"), 10, 64)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	var dto d.MovieGenresDto
	err = pkg.ReadJSON(w, r, &dto)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	// Call service method
	genres, err := h.MovieService.AddGenres(
		r.Context(),
		movieId,
		dto.Genres,
	)

	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	pkg.Ok(genres, "Genres added successfully", w)
}

// DeleteGenres removes specified genres from a movie
func (h *MoviesHandler) DeleteMovieGenres(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.ParseUint(chi.URLParam(r, "movieId"), 10, 64)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	var dto d.MovieGenresDto
	err = pkg.ReadJSON(w, r, &dto)
	if err != nil {
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}
	if errs := pkg.ValidateStruct(dto); errs != nil {
		pkg.BadRequest(w, r, &errs)
		return
	}

	// Call service method
	err = h.MovieService.RemoveGenres(
		r.Context(),
		movieId,
		dto.Genres,
	)

	if err != nil {
		if errors.Is(err, e.ErrMovieOneGenre) {
			pkg.Error(e.ErrMovieOneGenreAPI, w, r)
			return
		}
		if errors.Is(err, e.ErrSomeGenresNotFound) {
			pkg.Error(e.ErrSomeGenresNotFoundAPI, w, r)
			return
		}
		pkg.BadRequest(w, r, (*any)(nil))
		return
	}

	pkg.Ok((*any)(nil), "Genres deleted successfully", w)
}
