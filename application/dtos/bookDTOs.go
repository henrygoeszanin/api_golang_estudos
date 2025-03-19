package dtos

import (
	"time"

	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// BookResponseDTO representa os dados de livro que serão retornados nas respostas da API
type BookResponseDTO struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Available   int       `json:"available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BookCreateDTO representa os dados para criação de um livro
type BookCreateDTO struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Author      string `json:"author" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"min=1"`
}

// BookUpdateDTO representa os dados para atualização de um livro
type BookUpdateDTO struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=200"`
	Author      string `json:"author" binding:"omitempty,min=1,max=100"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"omitempty,min=1"`
}

// BookToResponseDTO converte uma entidade Book para um BookResponseDTO
func BookToResponseDTO(book entities.Book) BookResponseDTO {
	return BookResponseDTO{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Description: book.Description,
		Quantity:    book.Quantity,
		Available:   book.Available,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}
