package entities

import (
	"gorm.io/gorm"
)

// User representa um usuário do sistema de biblioteca
type User struct {
	gorm.Model
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;not null;unique"`
	Password string `gorm:"size:255;not null"`
	IsAdmin  bool   `gorm:"default:false"`
	Loans    []Loan
}
