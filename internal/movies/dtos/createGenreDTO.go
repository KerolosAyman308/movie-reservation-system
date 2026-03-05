package movies

type CreateGenreDTO struct {
	Name string `json:"name"  validate:"required,max=100,min=5"`
}
