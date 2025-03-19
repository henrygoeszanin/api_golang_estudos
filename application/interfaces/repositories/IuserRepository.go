package repositories

import (
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
)

// UserRepository define as operações possíveis no repositório de usuários
type UserRepository interface {
	Create(user *entities.User) error
	FindByID(id uint) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id uint) error
	List() ([]*entities.User, error)
	PromoteToAdmin(id uint) error
	IsFirstUser() (bool, error)
}
