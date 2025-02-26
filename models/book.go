package models

import "gorm.io/gorm"

// Book Model
type Book struct {
	gorm.Model
	Title      string `json:"title" gorm:"not null"`
	Author     string `json:"author" gorm:"not null"`
	CoverImage string `json:"cover_image"`
	Content    string `json:"content" gorm:"type:text"`
	Likes      int    `json:"likes" gorm:"default:0"`
}
