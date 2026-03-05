package movies

import (
	"context"
	"io"
	f "movie/system/internal/files"
	d "movie/system/internal/movies/dtos"
	m "movie/system/internal/movies/models"
	"movie/system/pkg"
)

// IMoviesService defines the contract for managing movie-related operations.
type IMoviesService interface {
	// GetPaginated retrieves a paginated list of movies.
	// It applies optional filtering by title, name, and genre (via subquery).
	// The response includes the total count of unique movies and eagerly loads associated genre data.
	GetPaginated(ctx context.Context, pag *pkg.PaginationRequest, searchByTitle string, searchByName string, searchByGenre string) (*pkg.PaginationResponse[m.Movie], error)

	// AddGenres links additional genres to an existing movie.
	// It safely skips any genres already attached to the movie to prevent duplicates.
	// Returns e.ErrMovieNotFound if the movie doesn't exist, or e.ErrSomeGenresNotFound if any provided genre ID is invalid.
	AddGenres(ctx context.Context, movieId uint64, genresId []uint64) ([]m.MovieGenres, error)

	// RemoveGenres detaches specific genres from a movie.
	// Business Rule: A movie must have at least one genre. If the operation attempts to remove
	// all existing genres, it will return an e.ErrMovieOneGenre error.
	// Also validates that all requested genre IDs are currently associated with the movie.
	RemoveGenres(ctx context.Context, movieId uint64, genresId []uint64) error

	// UploadImage uploads a file to the storage bucket and links it to the specified movie.
	// If the movie already has an image associated with it, the existing image is deleted from the bucket first.
	// If the database update fails, the newly uploaded file is safely cleaned up (deleted).
	UploadImage(ctx context.Context, movieId uint64, file io.Reader, fileName string) (*f.File, error)

	// DeleteImage removes the image file associated with a movie from the storage bucket.
	// Note: Returns e.ErrMovieNotFound if the movie record does not exist.
	DeleteImage(ctx context.Context, movieId uint64) error

	// Create inserts a newly created movie and its associated genres into the database.
	// It executes within a transaction to ensure data integrity.
	// Business Rules:
	// - Title must be unique case-insensitively (returns e.ErrMovieDubName if duplicate).
	// - All provided genre IDs must exist in the database (returns e.ErrSomeGenresNotFound otherwise).
	Create(ctx context.Context, dto *d.CreateMovieDTO) (*m.Movie, error)

	// Update modifies an existing movie's details.
	// Note: This method is currently unimplemented and will panic if called.
	Update(ctx context.Context, id uint64, dto *d.CreateMovieDTO) (*m.Movie, error)

	// Delete permanently removes a movie and its genre associations via a database transaction.
	// If the database deletion is successful, it also attempts to delete the associated image
	// from the storage bucket (logging an error if the file deletion fails, but still returning success).
	Delete(ctx context.Context, id uint64) error
}
