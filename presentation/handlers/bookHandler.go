package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/services"
)

// BookHandler manipula as requisições relacionadas a livros
type BookHandler struct {
	bookService services.BookService
}

// NewBookHandler cria uma nova instância de BookHandler
func NewBookHandler(bookService services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// List lista todos os livros
func (h *BookHandler) List(c *gin.Context) {
	books, err := h.bookService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetByID busca um livro pelo ID
func (h *BookHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	book, err := h.bookService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Create adiciona um novo livro
func (h *BookHandler) Create(c *gin.Context) {
	var bookDTO dtos.BookCreateDTO
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBook, err := h.bookService.Create(bookDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdBook)
}

// Update atualiza os dados de um livro
func (h *BookHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var bookDTO dtos.BookUpdateDTO
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.bookService.Update(uint(id), bookDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Delete remove um livro
func (h *BookHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.bookService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Livro removido com sucesso"})
}
