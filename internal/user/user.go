package user

import (
	log "log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint    `gorm:"autoIncrement;primaryKey"`
	Email     string  `gorm:"uniqueIndex;not null;type:varchar(255)"`
	FirstName string  `gorm:"not null;type:varchar(50)"`
	LastName  *string `gorm:"type:varchar(50)"`
	IsAdmin   bool    `gorm:"default:false"`
	Password  string  `gorm:"not null;"`
	IsActive  bool    `gorm:"default:true"`
	Birthday  time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *User) SetPassword() error {
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errHash != nil {
		log.Error("Error at hashing password", "Error:", errHash.Error())
		return errHash
	}
	u.Password = string(hashedPassword)

	return errHash
}

func (u *User) IsValidPassword(pass string) bool {
	isSuccess := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))

	return isSuccess != nil
}
