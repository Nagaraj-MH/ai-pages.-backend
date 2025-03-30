package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `gorm:"unique" json:"email"`
	Password     string `json:"-"`
	ResetToken   string `json:"reset_token"`
	TokenExpires time.Time
}
