package dtos

import (
	"time"

	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// UserCreateDTO representa os dados para criação de um usuário
type UserCreateDTO struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateDTO representa os dados para atualização de um usuário
type UserUpdateDTO struct {
	Name     string `json:"name" binding:"omitempty,min=3,max=100"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

// UserResponseDTO representa os dados de usuário que serão retornados nas respostas da API
type UserResponseDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponseDTO converte uma entidade User para um UserResponseDTO
func ToResponseDTO(user entities.User) UserResponseDTO {
	return UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
