package repositories

import (
	"errors"

	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
	"gorm.io/gorm"
)

// userRepository implementa a interface UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository cria uma nova instância do repositório de usuários
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create cria um novo usuário no banco de dados
func (r *userRepository) Create(user *entities.User) error {
	// Verificar se é o primeiro usuário
	isFirst, err := r.IsFirstUser()
	if err != nil {
		return err
	}

	// Se for o primeiro usuário, definir como administrador
	if isFirst {
		user.IsAdmin = true
	}

	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByID busca um usuário pelo seu ID
func (r *userRepository) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Usuário não encontrado
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail busca um usuário pelo seu email
func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Usuário não encontrado
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update atualiza os dados de um usuário
func (r *userRepository) Update(user *entities.User) error {
	result := r.db.Save(user)
	return result.Error
}

// Delete remove um usuário pelo seu ID
func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&entities.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("usuário não encontrado")
	}
	return nil
}

// List retorna todos os usuários
func (r *userRepository) List() ([]*entities.User, error) {
	var users []*entities.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// PromoteToAdmin promove um usuário para administrador
func (r *userRepository) PromoteToAdmin(id uint) error {
	result := r.db.Model(&entities.User{}).Where("id = ?", id).Update("is_admin", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("usuário não encontrado")
	}
	return nil
}

// IsFirstUser verifica se este será o primeiro usuário no sistema
func (r *userRepository) IsFirstUser() (bool, error) {
	var count int64
	if err := r.db.Model(&entities.User{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
