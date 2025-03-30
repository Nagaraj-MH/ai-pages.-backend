package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title      string `json:"title" gorm:"not null"`
	Author     string `json:"author" gorm:"not null"`
	CoverImage string `json:"cover_image"`
	PDFData    []byte `json:"pdf_data" gorm:"type:bytea"` 
	Likes      int    `json:"likes" gorm:"default:0"`
}

