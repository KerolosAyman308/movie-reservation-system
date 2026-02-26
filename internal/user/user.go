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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password", "error", err)
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// IsValidPassword returns true if the provided plain-text password matches the stored hash.
// bcrypt.CompareHashAndPassword returns nil on a successful match.
func (u *User) IsValidPassword(pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	return err == nil
}
