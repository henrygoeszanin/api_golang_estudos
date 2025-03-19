package repositories

import (
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// BookRepository define as operações possíveis no repositório de livros
type BookRepository interface {
	Create(book *entities.Book) error
	FindByID(id uint) (*entities.Book, error)
	List() ([]*entities.Book, error)
	Update(book *entities.Book) error
	Delete(id uint) error
}
