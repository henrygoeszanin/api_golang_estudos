package entities

import (
	"time"

	"gorm.io/gorm"
)

// Loan representa um empréstimo de livro
type Loan struct {
	gorm.Model
	UserID     uint
	User       User `gorm:"foreignKey:UserID"`
	BookID     uint
	Book       Book      `gorm:"foreignKey:BookID"`
	LoanDate   time.Time `gorm:"not null"`
	ReturnDate time.Time `gorm:"not null"` // Data prevista para devolução
	ReturnedAt *time.Time
	IsReturned bool `gorm:"default:false"`
}
