package entities

import (
	"gorm.io/gorm"
)

// Book representa um livro na biblioteca
type Book struct {
	gorm.Model
	Title       string `gorm:"size:200;not null"`
	Author      string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
	Quantity    int    `gorm:"default:1"`
	Available   int    `gorm:"default:1"`
	Loans       []Loan
}
