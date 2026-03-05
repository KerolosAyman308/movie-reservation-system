package movies

import (
	"context"
	d "movie/system/internal/movies/dtos"
	m "movie/system/internal/movies/models"
	"movie/system/pkg"
)

// IGenreService defines the contract for managing movie genres.
type IGenreService interface {
	// GetPaginated retrieves a paginated list of genres based on the provided pagination request.
	// It supports optional dynamic sorting (by SortKey and Sort direction) and includes
	// the total count of all genres in the database within the response.
	GetPaginated(ctx context.Context, pag *pkg.PaginationRequest) (*pkg.PaginationResponse[m.Genre], error)

	// Create inserts a new genre into the database.
	// Business Rule: Genre names must be strictly unique. The implementation performs a
	// case-insensitive and space-trimmed check before creation.
	// Returns e.ErrGenreDubName if a genre with the same name already exists.
	Create(ctx context.Context, dto *d.CreateGenreDTO) (*m.Genre, error)

	// Delete removes a genre from the database by its ID.
	// Business Rules:
	// - Returns e.ErrGenreNotFound if the target ID does not exist.
	// - Prevents deletion if the genre is currently associated with any movies,
	//   returning e.ErrGenreHasMovies to enforce referential integrity.
	Delete(ctx context.Context, id uint64) error

	// FindByIds retrieves a list of genres that match the provided slice of IDs.
	// Note: This method returns only the genres that were successfully found.
	// The caller is responsible for verifying if the length of the returned slice
	// matches the length of the requested IDs if strict existence validation is required.
	FindByIds(ctx context.Context, ids []uint64) ([]m.Genre, error)
}
