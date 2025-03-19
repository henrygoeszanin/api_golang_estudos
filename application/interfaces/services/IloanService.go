package services

import (
	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
)

// LoanService define os serviços disponíveis para empréstimos
type LoanService interface {
	Create(userID uint, loanDTO dtos.LoanCreateDTO) (*dtos.LoanResponseDTO, error)
	GetByID(id uint, userID uint) (*dtos.LoanResponseDTO, error)
	ListByUser(userID uint) ([]dtos.LoanResponseDTO, error)
	ReturnLoan(id uint, userID uint) (*dtos.LoanResponseDTO, error)
}
