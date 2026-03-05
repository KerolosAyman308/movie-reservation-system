package movies

type Genre struct {
	Id uint64 `gorm:"autoIncrement;primaryKey"`

	Name string `gorm:"uniqueIndex;not null;type: varchar(100)"`

	Movies []MovieGenres
}
