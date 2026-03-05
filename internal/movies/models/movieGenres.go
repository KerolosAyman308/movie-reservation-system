package movies

type MovieGenres struct {
	Id uint64 `gorm:"autoIncrement;primaryKey"`

	MovieId uint64 `gorm:"uniqueIndex:idx_movie_genre"`
	Movie   Movie  `gorm:"foreignkey:MovieId"`

	GenreId uint64 `gorm:"uniqueIndex:idx_movie_genre"`
	Genre   Genre  `gorm:"foreignkey:GenreId"`
}
