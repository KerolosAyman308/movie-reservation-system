package movies

type CreateMovieDTO struct {
	Title string `json:"title" validate:"required,max=100,min=5"`

	Description string `json:"description" validate:"required,max=350,min=5"`

	Genres []uint64 `json:"genres" validate:"required,min=1,dive"`
}

type UpdateMovieDTO struct {
	Description string `json:"description" validate:"required,max=350,min=5"`

	Genres []uint64 `json:"genres" validate:"required,min=1,dive"`
}

type MovieGenresDto struct {
	Genres []uint64 `json:"genres"  validate:"required,dive,min=1"`
}
