package repositories

import (
	"errors"

	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
	"gorm.io/gorm"
)

// bookRepository implementa a interface BookRepository
type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository cria uma nova instância do repositório de livros
func NewBookRepository(db *gorm.DB) repositories.BookRepository {
	return &bookRepository{
		db: db,
	}
}

// Create cria um novo livro no banco de dados
func (bookRepository *bookRepository) Create(book *entities.Book) error {
	// Garantir que disponível = quantidade inicialmente
	book.Available = book.Quantity

	result := bookRepository.db.Create(book)
	return result.Error
}

// FindByID busca um livro pelo seu ID
func (bookRepository *bookRepository) FindByID(id uint) (*entities.Book, error) {
	var book entities.Book
	result := bookRepository.db.First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Livro não encontrado
		}
		return nil, result.Error
	}
	return &book, nil
}

// List retorna todos os livros
func (bookRepository *bookRepository) List() ([]*entities.Book, error) {
	var books []*entities.Book
	result := bookRepository.db.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

// Update atualiza os dados de um livro
func (bookRepository *bookRepository) Update(book *entities.Book) error {
	result := bookRepository.db.Save(book)
	return result.Error
}

// Delete remove um livro pelo seu ID
func (bookRepository *bookRepository) Delete(id uint) error {
	result := bookRepository.db.Delete(&entities.Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("livro não encontrado")
	}
	return nil
}
