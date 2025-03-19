package dtos

import (
	"time"

	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// LoanCreateDTO representa os dados para criação de um empréstimo
type LoanCreateDTO struct {
	BookID     uint      `json:"book_id" binding:"required"`
	ReturnDate time.Time `json:"return_date" binding:"required,gt"`
}

// LoanResponseDTO representa os dados de empréstimo que serão retornados nas respostas da API
type LoanResponseDTO struct {
	ID         uint       `json:"id"`
	BookID     uint       `json:"book_id"`
	BookTitle  string     `json:"book_title"`
	UserID     uint       `json:"user_id"`
	UserName   string     `json:"user_name"`
	LoanDate   time.Time  `json:"loan_date"`
	ReturnDate time.Time  `json:"return_date"`
	ReturnedAt *time.Time `json:"returned_at"`
	IsReturned bool       `json:"is_returned"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// LoanToResponseDTO converte uma entidade Loan para um LoanResponseDTO
func LoanToResponseDTO(loan entities.Loan) LoanResponseDTO {
	return LoanResponseDTO{
		ID:         loan.ID,
		BookID:     loan.BookID,
		BookTitle:  loan.Book.Title,
		UserID:     loan.UserID,
		UserName:   loan.User.Name,
		LoanDate:   loan.LoanDate,
		ReturnDate: loan.ReturnDate,
		ReturnedAt: loan.ReturnedAt,
		IsReturned: loan.IsReturned,
		CreatedAt:  loan.CreatedAt,
		UpdatedAt:  loan.UpdatedAt,
	}
}
