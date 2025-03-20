package services

import (
	"errors"

	"github.com/henrygoeszanin/api_golang_estudos/application/dtos"
	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/repositories"
	"github.com/henrygoeszanin/api_golang_estudos/application/interfaces/services"
	"github.com/henrygoeszanin/api_golang_estudos/domain/entities"
	"golang.org/x/crypto/bcrypt"
)

// userService implementa a interface UserService
type userService struct {
	userRepository repositories.UserRepository
}

// NewUserService cria uma nova instância do serviço de usuários
func NewUserService(userRepository repositories.UserRepository) services.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

// Create cria um novo usuário
func (userService *userService) Create(userDTO dtos.UserCreateDTO) (*dtos.UserResponseDTO, error) {
	// Verificar se o email já está em uso
	existingUser, err := userService.userRepository.FindByEmail(userDTO.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email já está em uso")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Criar novo usuário
	user := entities.User{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: string(hashedPassword),
	}

	// Salvar no banco de dados
	if err := userService.userRepository.Create(&user); err != nil {
		return nil, err
	}

	// Converter para o DTO de resposta
	responseDTO := dtos.ToResponseDTO(user)

	return &responseDTO, nil
}

// GetByID busca um usuário pelo ID
func (userService *userService) GetByID(id uint) (*dtos.UserResponseDTO, error) {
	user, err := userService.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	responseDTO := dtos.ToResponseDTO(*user)
	return &responseDTO, nil
}

// GetByEmail busca um usuário pelo email
func (userService *userService) GetByEmail(email string) (*entities.User, error) {
	return userService.userRepository.FindByEmail(email)
}

// Update atualiza os dados de um usuário
func (userService *userService) Update(id uint, userDTO dtos.UserUpdateDTO) (*dtos.UserResponseDTO, error) {
	user, err := userService.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Atualizar campos se fornecidos
	if userDTO.Name != "" {
		user.Name = userDTO.Name
	}
	if userDTO.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	// Salvar alterações
	if err := userService.userRepository.Update(user); err != nil {
		return nil, err
	}

	responseDTO := dtos.ToResponseDTO(*user)
	return &responseDTO, nil
}

// Delete remove um usuário
func (userService *userService) Delete(id uint) error {
	return userService.userRepository.Delete(id)
}

// List retorna todos os usuários
func (userService *userService) List() ([]dtos.UserResponseDTO, error) {
	users, err := userService.userRepository.List()
	if err != nil {
		return nil, err
	}

	var userDTOs []dtos.UserResponseDTO
	for _, user := range users {
		userDTOs = append(userDTOs, dtos.ToResponseDTO(*user))
	}

	return userDTOs, nil
}

// PromoteToAdmin promove um usuário para administrador
func (userService *userService) PromoteToAdmin(id uint) (*dtos.UserResponseDTO, error) {
	if err := userService.userRepository.PromoteToAdmin(id); err != nil {
		return nil, err
	}

	user, err := userService.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	responseDTO := dtos.ToResponseDTO(*user)
	return &responseDTO, nil
}

// AuthenticateUser autentica um usuário pelo email e senha
func (userService *userService) AuthenticateUser(email, password string) (*entities.User, error) {
	user, err := userService.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	return user, nil
}
