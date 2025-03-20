package services

import (
	"errors"
	"time"

	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/services"
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// loanService implementa a interface LoanService
type loanService struct {
	loanRepository repositories.LoanRepository
	bookRepository repositories.BookRepository
}

// NewLoanService cria uma nova instância do serviço de empréstimos
func NewLoanService(loanRepository repositories.LoanRepository, bookRepository repositories.BookRepository) services.LoanService {
	return &loanService{
		loanRepository: loanRepository,
		bookRepository: bookRepository,
	}
}

// Create cria um novo empréstimo
func (loanService *loanService) Create(userID uint, loanDTO dtos.LoanCreateDTO) (*dtos.LoanResponseDTO, error) {
	// Verificar se o livro existe
	book, err := loanService.bookRepository.FindByID(loanDTO.BookID)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("livro não encontrado")
	}

	// Verificar se há exemplares disponíveis
	if book.Available <= 0 {
		return nil, errors.New("livro não disponível para empréstimo")
	}

	// Criar empréstimo
	loan := entities.Loan{
		UserID:     userID,
		BookID:     loanDTO.BookID,
		LoanDate:   time.Now(),
		ReturnDate: loanDTO.ReturnDate,
		IsReturned: false,
	}

	if err := loanService.loanRepository.Create(&loan); err != nil {
		return nil, err
	}

	// Carregar dados completos do empréstimo com livro e usuário
	fullLoan, err := loanService.loanRepository.FindByID(loan.ID)
	if err != nil {
		return nil, err
	}

	responseDTO := dtos.LoanToResponseDTO(*fullLoan)
	return &responseDTO, nil
}

// GetByID busca um empréstimo pelo ID
func (loanService *loanService) GetByID(id uint, userID uint) (*dtos.LoanResponseDTO, error) {
	loan, err := loanService.loanRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, errors.New("empréstimo não encontrado")
	}

	// Verificar se o empréstimo pertence ao usuário
	if loan.UserID != userID {
		return nil, errors.New("acesso negado a este empréstimo")
	}

	responseDTO := dtos.LoanToResponseDTO(*loan)
	return &responseDTO, nil
}

// ListByUser retorna todos os empréstimos de um usuário
func (loanService *loanService) ListByUser(userID uint) ([]dtos.LoanResponseDTO, error) {
	loans, err := loanService.loanRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var loanDTOs []dtos.LoanResponseDTO
	for _, loan := range loans {
		loanDTOs = append(loanDTOs, dtos.LoanToResponseDTO(*loan))
	}

	return loanDTOs, nil
}

// ReturnLoan marca um empréstimo como devolvido
func (loanService *loanService) ReturnLoan(id uint, userID uint) (*dtos.LoanResponseDTO, error) {
	// Verificar se o empréstimo existe e pertence ao usuário
	loan, err := loanService.loanRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, errors.New("empréstimo não encontrado")
	}

	if loan.UserID != userID {
		return nil, errors.New("acesso negado a este empréstimo")
	}

	if loan.IsReturned {
		return nil, errors.New("empréstimo já foi devolvido")
	}

	// Processar devolução
	returnDate := time.Now()
	if err := loanService.loanRepository.ReturnLoan(id, returnDate); err != nil {
		return nil, err
	}

	// Obter empréstimo atualizado
	updatedLoan, err := loanService.loanRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	responseDTO := dtos.LoanToResponseDTO(*updatedLoan)
	return &responseDTO, nil
}
