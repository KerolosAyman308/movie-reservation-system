package movies

import (
	"context"
	"errors"
	"fmt"
	"io"
	logs "log/slog"
	f "movie/system/internal/files"
	e "movie/system/internal/movies"
	d "movie/system/internal/movies/dtos"
	mo "movie/system/internal/movies/models"
	"movie/system/pkg"
	"strings"

	"gorm.io/gorm"
)

type MoviesService struct {
	DB           *gorm.DB
	GenreService IGenreService
	FileService  f.IFiles
	BucketName   string
}

func NewMoviesService(db *gorm.DB, genreService IGenreService, fileService f.IFiles, bucketName string) IMoviesService {
	return &MoviesService{
		DB:           db,
		GenreService: genreService,
		FileService:  fileService,
		BucketName:   bucketName,
	}
}

// Create implements IMoviesService.
func (m *MoviesService) Create(ctx context.Context, dto *d.CreateMovieDTO) (*mo.Movie, error) {
	//check if title exists
	var count int64
	err := m.DB.Model(&mo.Movie{}).Where("LOWER(title) = ?", strings.TrimSpace(strings.ToLower(dto.Title))).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, e.ErrMovieDubName
	}

	// Validate the genres are already exists
	genres, err := m.GenreService.FindByIds(ctx, dto.Genres)
	if err != nil {
		return nil, err
	}
	if len(genres) != len(dto.Genres) {
		return nil, e.ErrSomeGenresNotFound
	}

	var createdMovie mo.Movie = mo.Movie{
		Title:       dto.Title,
		Description: dto.Description,
	}

	tx := m.DB.Begin()

	defer func() {
		tx.Rollback()
	}()

	if err := tx.Model(&mo.Movie{}).WithContext(ctx).Create(&createdMovie).Error; err != nil {
		return nil, err
	}

	var movieGenres []mo.MovieGenres = make([]mo.MovieGenres, len(dto.Genres))
	for i, genre := range genres {
		movieGenres[i] = mo.MovieGenres{
			GenreId: genre.Id,
			MovieId: createdMovie.Id,
			Genre:   genre, // I put this here so it will be easy in populating and sending to the user without needing for extra query
		}
	}

	if err := tx.WithContext(ctx).Model(&mo.MovieGenres{}).Create(&movieGenres).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	createdMovie.Genres = movieGenres
	return &createdMovie, nil
}

// AddGenres implements IMoviesService.
func (m *MoviesService) AddGenres(ctx context.Context, movieId uint64, genresId []uint64) ([]mo.MovieGenres, error) {
	//check movie exist
	var currMovieGenres []mo.MovieGenres
	err := m.DB.Model(&mo.MovieGenres{}).WithContext(ctx).Where("movie_id = ?", movieId).Find(&currMovieGenres).Error
	if err != nil {
		return nil, err
	}
	if len(currMovieGenres) == 0 {
		return nil, e.ErrMovieNotFound
	}
	//check genre exist
	genres, err := m.GenreService.FindByIds(ctx, genresId)
	if err != nil {
		return nil, err
	}
	if len(genres) != len(genresId) {
		return nil, e.ErrSomeGenresNotFound
	}

	//insert movie genre
	var movieGenres []mo.MovieGenres = []mo.MovieGenres{}
	currMovieGenresMap := make(map[int]bool)
	for _, val := range currMovieGenres {
		currMovieGenresMap[int(val.GenreId)] = true
	}
	var i int = 0
	for _, genre := range genresId {
		_, existed := currMovieGenresMap[int(genre)]
		if existed {
			continue
		}
		movieGenres = append(movieGenres, mo.MovieGenres{
			GenreId: genre,
			MovieId: movieId,
		})
		i++
	}

	if len(movieGenres) == 0 {
		return movieGenres, nil
	}
	if err := m.DB.Model(&mo.MovieGenres{}).WithContext(ctx).Create(&movieGenres).Error; err != nil {
		return nil, err
	}

	return movieGenres, nil
}

// RemoveGenres implements IMoviesService.
func (m *MoviesService) RemoveGenres(ctx context.Context, movieId uint64, movieGenreIds []uint64) error {
	var movieGenres []mo.MovieGenres
	err := m.DB.Model(&mo.MovieGenres{}).WithContext(ctx).Where("movie_id = ?", movieId).Find(&movieGenres).Error
	if err != nil {
		return err
	}
	// to keep at least one genre attached to the movie
	if (len(movieGenres) - 1) < len(movieGenreIds) {
		return e.ErrMovieOneGenre
	}

	// Check if all genres exist
	movieGenresMap := make(map[int]uint64)
	for _, value := range movieGenres {
		movieGenresMap[int(value.GenreId)] = value.Id
	}

	movieGenresToDelete := []uint64{}
	for _, value := range movieGenreIds {
		v, exists := movieGenresMap[int(value)]
		if !exists {
			return e.ErrSomeGenresNotFound
		}
		movieGenresToDelete = append(movieGenresToDelete, v)
	}

	err = m.DB.Model(&mo.MovieGenres{}).WithContext(ctx).Delete("id IN ?", movieGenresToDelete).Error
	return err
}

// DeleteImage implements IMoviesService.
func (m *MoviesService) DeleteImage(ctx context.Context, movieId uint64) error {
	var movie mo.Movie
	err := m.DB.Model(&mo.Movie{}).WithContext(ctx).Preload("Image").First(movieId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrMovieNotFound
		}
		return err
	}

	if movie.Image != nil {
		if err := m.FileService.Delete(ctx, movie.Image.BucketName, *movie.ImageObjectKey); err != nil {
			return err
		}
	}

	return nil
}

// UploadImage implements IMoviesService.
func (m *MoviesService) UploadImage(ctx context.Context, movieId uint64, file io.Reader, fileName string) (*f.File, error) {
	var movie mo.Movie
	err := m.DB.Model(&mo.Movie{}).WithContext(ctx).Preload("Image").First(&movie, movieId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrMovieNotFound
		}
		return nil, err
	}

	if movie.Image != nil {
		if err := m.FileService.Delete(ctx, movie.Image.BucketName, *movie.ImageObjectKey); err != nil {
			return nil, err
		}
	}

	createdFile, err := m.FileService.Create(ctx, m.BucketName, fileName, file)
	if err != nil {
		return nil, err
	}
	err = m.DB.Model(&mo.Movie{}).WithContext(ctx).Where("id = ?", movieId).UpdateColumn("image_object_key", createdFile.ObjectKey).Error
	if err != nil {
		m.FileService.Delete(ctx, createdFile.BucketName, createdFile.ObjectKey)
		return nil, err
	}
	return createdFile, nil
}

// Delete implements IMoviesService.
func (m *MoviesService) Delete(ctx context.Context, id uint64) error {
	var movie mo.Movie
	err := m.DB.Model(&mo.Movie{}).WithContext(ctx).Preload("Image").First(&movie, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrMovieNotFound
		}
		return err
	}

	err = m.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("movie_id = ?", movie.Id).Delete(&mo.MovieGenres{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&movie).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if movie.Image != nil {
		if err := m.FileService.Delete(ctx, movie.Image.BucketName, *movie.ImageObjectKey); err != nil {
			logs.Error(err.Error())
		}
	}
	return nil
}

func (m *MoviesService) GetPaginated(ctx context.Context, pag *pkg.PaginationRequest, searchByTitle string, searchByName string, searchByGenre string) (*pkg.PaginationResponse[mo.Movie], error) {
	var movies []mo.Movie
	var size int64

	search := m.DB.WithContext(ctx).Model(&mo.Movie{})

	if searchByTitle != "" {
		search = search.Where("title LIKE ?", "%"+searchByTitle+"%")
	}
	if searchByName != "" {
		search = search.Where("name LIKE ?", "%"+searchByName+"%")
	}

	if searchByGenre != "" {
		subQuery := m.DB.Model(&mo.MovieGenres{}).
			Select("movie_genres.movie_id").
			Joins("JOIN genres ON genres.id = movie_genres.genre_id").
			Where("genres.name LIKE ?", "%"+searchByGenre+"%")

		search = search.Where("id IN (?)", subQuery)
	}

	// 3. Get the accurate total count of UNIQUE movies
	if err := search.Count(&size).Error; err != nil {
		return nil, err
	}

	if pag.SortKey != "" && pag.Sort != "" {
		search = search.Order(fmt.Sprintf("%s %s", pag.SortKey, pag.Sort))
	}

	err := search.
		Limit(pag.Limit).
		Offset(pag.GetOffset()).
		Preload("Image").
		Preload("Genres").
		Preload("Genres.Genre").
		Find(&movies).Error

	if err != nil {
		return nil, err
	}

	for _, movie := range movies {
		if movie.Image != nil {
			movie.Image.Url, _ = m.FileService.GenerateUrl(ctx, movie.Image.BucketName, movie.Image.ObjectKey)
		}
	}

	return &pkg.PaginationResponse[mo.Movie]{
		Data: movies,
		Size: size,
	}, nil
}

// Update implements IMoviesService.
func (m *MoviesService) Update(ctx context.Context, id uint64, dto *d.CreateMovieDTO) (*mo.Movie, error) {
	panic("unimplemented")
}
