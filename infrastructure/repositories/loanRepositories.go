package repositories

import (
	"errors"
	"time"

	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
	"gorm.io/gorm"
)

// loanRepository implementa a interface LoanRepository
type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository cria uma nova instância do repositório de empréstimos
func NewLoanRepository(db *gorm.DB) repositories.LoanRepository {
	return &loanRepository{
		db: db,
	}
}

// Create cria um novo empréstimo no banco de dados
func (loanRepository *loanRepository) Create(loan *entities.Loan) error {
	// Verificar se o livro está disponível
	var book entities.Book
	if err := loanRepository.db.First(&book, loan.BookID).Error; err != nil {
		return err
	}

	if book.Available <= 0 {
		return errors.New("livro não disponível para empréstimo")
	}

	// Iniciar transação
	tx := loanRepository.db.Begin()

	// Diminuir contador de disponíveis
	if err := tx.Model(&book).Update("available", book.Available-1).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Criar empréstimo
	if err := tx.Create(loan).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// FindByID busca um empréstimo pelo seu ID
func (loanRepository *loanRepository) FindByID(id uint) (*entities.Loan, error) {
	var loan entities.Loan
	result := loanRepository.db.Preload("Book").Preload("User").First(&loan, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Empréstimo não encontrado
		}
		return nil, result.Error
	}
	return &loan, nil
}

// FindByUserID busca todos os empréstimos de um usuário
func (loanRepository *loanRepository) FindByUserID(userID uint) ([]*entities.Loan, error) {
	var loans []*entities.Loan
	result := loanRepository.db.Where("user_id = ?", userID).Preload("Book").Preload("User").Find(&loans)
	if result.Error != nil {
		return nil, result.Error
	}
	return loans, nil
}

// Update atualiza os dados de um empréstimo
func (loanRepository *loanRepository) Update(loan *entities.Loan) error {
	result := loanRepository.db.Save(loan)
	return result.Error
}

// ReturnLoan marca um empréstimo como devolvido e atualiza o estoque do livro
func (loanRepository *loanRepository) ReturnLoan(id uint, returnDate time.Time) error {
	// Obter o empréstimo
	var loan entities.Loan
	if err := loanRepository.db.First(&loan, id).Error; err != nil {
		return err
	}

	if loan.IsReturned {
		return errors.New("empréstimo já foi devolvido")
	}

	// Iniciar transação
	tx := loanRepository.db.Begin()

	// Atualizar empréstimo
	if err := tx.Model(&loan).Updates(map[string]interface{}{
		"is_returned": true,
		"returned_at": returnDate,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Aumentar disponibilidade do livro
	if err := tx.Model(&entities.Book{}).Where("id = ?", loan.BookID).
		Update("available", gorm.Expr("available + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
