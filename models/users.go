package models

import (
	"megrez/libs/crypto"
	"megrez/services/config"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID uint `json:"id" gorm:"primary_key;autoIncrement;index"`

	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex;unique;not null"`
	Password string `json:"password,omitempty" gorm:"type:varchar(255);not null"`
	Role     int    `json:"role" gorm:"not null,default:0"`

	Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;unique;not null"`

	Balance float64 `json:"balance" gorm:"not null"`

	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Users) PasswordHash(password string) string {
	return crypto.Sha256(password + u.Email + config.GetSystemSalt())
}

func (u *Users) CheckPassword(password string) bool {
	return u.Password == u.PasswordHash(password)
}
