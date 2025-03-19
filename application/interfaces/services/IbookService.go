package services

import (
	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
)

// BookService define os serviços disponíveis para livros
type BookService interface {
	Create(bookDTO dtos.BookCreateDTO) (*dtos.BookResponseDTO, error)
	GetByID(id uint) (*dtos.BookResponseDTO, error)
	List() ([]dtos.BookResponseDTO, error)
	Update(id uint, bookDTO dtos.BookUpdateDTO) (*dtos.BookResponseDTO, error)
	Delete(id uint) error
}
