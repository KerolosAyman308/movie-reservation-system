package movies

import (
	"context"
	"errors"
	"fmt"
	e "movie/system/internal/movies"
	d "movie/system/internal/movies/dtos"
	m "movie/system/internal/movies/models"
	"movie/system/pkg"
	"strings"

	"gorm.io/gorm"
)

type GenreService struct {
	DB *gorm.DB
}

func NewGenreService(db *gorm.DB) IGenreService {
	return &GenreService{
		DB: db,
	}
}

func (g *GenreService) GetPaginated(ctx context.Context, pag *pkg.PaginationRequest) (*pkg.PaginationResponse[m.Genre], error) {
	var genres []m.Genre
	query := g.DB.WithContext(ctx).Model(&m.Genre{}).Limit(pag.Limit).Offset(pag.GetOffset())
	if pag.Sort != "" && pag.SortKey != "" {
		query.Order(fmt.Sprintf("%s %s", pag.SortKey, pag.Sort))
	}
	err := query.Find(&genres).Error
	if err != nil {
		return nil, err
	}
	var genresLength int64
	err = g.DB.WithContext(ctx).Model(&m.Genre{}).Count(&genresLength).Error
	if err != nil {
		return nil, err
	}

	return &pkg.PaginationResponse[m.Genre]{
		Data: genres,
		Size: genresLength,
	}, nil
}

// FindByIds implements IGenreService.
func (g *GenreService) FindByIds(ctx context.Context, ids []uint64) ([]m.Genre, error) {
	var genres []m.Genre
	err := g.DB.Model(&m.Genre{}).WithContext(ctx).Where("Id IN ?", ids).Find(&genres).Error
	if err != nil {
		return nil, err
	}
	return genres, nil
}

func (g *GenreService) Create(ctx context.Context, dto *d.CreateGenreDTO) (*m.Genre, error) {
	var genreWithSameName int64
	err := g.DB.WithContext(ctx).Model(&m.Genre{}).Where("LOWER(name) = ?", strings.TrimSpace(strings.ToLower(dto.Name))).Count(&genreWithSameName).Error
	if err != nil {
		return nil, err
	}
	if genreWithSameName > 0 {
		return nil, e.ErrGenreDubName
	}

	var genre m.Genre = m.Genre{
		Name: dto.Name,
	}
	err = g.DB.WithContext(ctx).Model(&genre).Create(&genre).Error
	if err != nil {
		return nil, err
	}

	return &genre, nil
}

func (g *GenreService) Delete(ctx context.Context, id uint64) error {
	var genreExistsWithMovies m.Genre
	err := g.DB.WithContext(ctx).Preload("Movies").Where("Id = ?", id).First(&genreExistsWithMovies).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrGenreNotFound
		}
		return err
	}

	if len(genreExistsWithMovies.Movies) > 0 {
		return e.ErrGenreHasMovies
	}

	err = g.DB.WithContext(ctx).Delete(&m.Genre{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}
