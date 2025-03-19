package services

import (
	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// UserService define os serviços disponíveis para usuários
type UserService interface {
	Create(userDTO dtos.UserCreateDTO) (*dtos.UserResponseDTO, error)
	GetByID(id uint) (*dtos.UserResponseDTO, error)
	GetByEmail(email string) (*entities.User, error)
	Update(id uint, userDTO dtos.UserUpdateDTO) (*dtos.UserResponseDTO, error)
	Delete(id uint) error
	List() ([]dtos.UserResponseDTO, error)
	PromoteToAdmin(id uint) (*dtos.UserResponseDTO, error)
	AuthenticateUser(email, password string) (*entities.User, error)
}
