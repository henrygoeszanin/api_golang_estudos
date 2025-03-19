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
func (r *loanRepository) Create(loan *entities.Loan) error {
	// Verificar se o livro está disponível
	var book entities.Book
	if err := r.db.First(&book, loan.BookID).Error; err != nil {
		return err
	}

	if book.Available <= 0 {
		return errors.New("livro não disponível para empréstimo")
	}

	// Iniciar transação
	tx := r.db.Begin()

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
func (r *loanRepository) FindByID(id uint) (*entities.Loan, error) {
	var loan entities.Loan
	result := r.db.Preload("Book").Preload("User").First(&loan, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Empréstimo não encontrado
		}
		return nil, result.Error
	}
	return &loan, nil
}

// FindByUserID busca todos os empréstimos de um usuário
func (r *loanRepository) FindByUserID(userID uint) ([]*entities.Loan, error) {
	var loans []*entities.Loan
	result := r.db.Where("user_id = ?", userID).Preload("Book").Preload("User").Find(&loans)
	if result.Error != nil {
		return nil, result.Error
	}
	return loans, nil
}

// Update atualiza os dados de um empréstimo
func (r *loanRepository) Update(loan *entities.Loan) error {
	result := r.db.Save(loan)
	return result.Error
}

// ReturnLoan marca um empréstimo como devolvido e atualiza o estoque do livro
func (r *loanRepository) ReturnLoan(id uint, returnDate time.Time) error {
	// Obter o empréstimo
	var loan entities.Loan
	if err := r.db.First(&loan, id).Error; err != nil {
		return err
	}

	if loan.IsReturned {
		return errors.New("empréstimo já foi devolvido")
	}

	// Iniciar transação
	tx := r.db.Begin()

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
