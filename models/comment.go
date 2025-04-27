package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	BookID  uint   `json:"book_id" gorm:"not null"`
	UserID  uint   `json:"userId"`
	Content string `json:"content" gorm:"type:text;not null"`
}
