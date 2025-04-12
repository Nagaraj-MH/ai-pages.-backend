package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `gorm:"unique" json:"email"`
	Password     string `json:"password"`
	ResetToken   string `json:"reset_token"`
	TokenExpires time.Time
}
