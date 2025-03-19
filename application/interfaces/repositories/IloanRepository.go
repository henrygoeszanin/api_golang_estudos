package repositories

import (
	"time"

	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// LoanRepository define as operações possíveis no repositório de empréstimos
type LoanRepository interface {
	Create(loan *entities.Loan) error
	FindByID(id uint) (*entities.Loan, error)
	FindByUserID(userID uint) ([]*entities.Loan, error)
	Update(loan *entities.Loan) error
	ReturnLoan(id uint, returnDate time.Time) error
}
