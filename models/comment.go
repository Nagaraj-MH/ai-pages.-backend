package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	BookID  uint   `json:"book_id" gorm:"not null"`
	UserID  uint   `json:"user_id"`
	Content string `json:"content" gorm:"type:text;not null"`
}
