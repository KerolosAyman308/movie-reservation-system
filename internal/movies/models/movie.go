package movies

import (
	f "movie/system/internal/files"
)

type Movie struct {
	Id uint64 `gorm:"autoIncrement;primaryKey"`

	Title string `gorm:"uniqueIndex;not null;type: varchar(100)"`

	Description string `gorm:"not null;type: varchar(500)"`

	ImageObjectKey *string `gorm:"type:varchar(40)"`
	Image          *f.File `gorm:"foreignKey:ImageObjectKey;references:ObjectKey;constraint:OnDelete:SET NULL;"`
	Genres         []MovieGenres
}
