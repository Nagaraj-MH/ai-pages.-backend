package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	BookID  uint   `json:"book_id" gorm:"not null"`
	Username  string `json:"username" gorm:"not null"`
	Content string `json:"content" gorm:"type:text;not null"`
}
