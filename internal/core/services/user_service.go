package services

import (
	"fmt"

	"github.com/danielgalindoj/auth-service/internal/core/domain/model"
	"github.com/danielgalindoj/auth-service/internal/core/ports"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetByID(id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Limpiar password hash
	user.PasswordHash = nil
	return user, nil
}

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Limpiar password hash
	user.PasswordHash = nil
	return user, nil
}

func (s *UserService) List(page, limit int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := s.userRepo.List(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error obteniendo usuarios: %w", err)
	}

	total, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("error contando usuarios: %w", err)
	}

	return users, total, nil
}

func (s *UserService) Update(userID uuid.UUID, updates map[string]interface{}) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Aplicar actualizaciones
	for field, value := range updates {
		switch field {
		case "first_name":
			if v, ok := value.(string); ok {
				user.FirstName = &v
			}
		case "last_name":
			if v, ok := value.(string); ok {
				user.LastName = &v
			}
		case "username":
			if v, ok := value.(string); ok {
				user.Username = &v
			}
		case "avatar_url":
			if v, ok := value.(string); ok {
				user.AvatarURL = &v
			}
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("error actualizando usuario: %w", err)
	}

	// Limpiar password hash
	user.PasswordHash = nil
	return user, nil
}

func (s *UserService) Delete(userID uuid.UUID) error {
	return s.userRepo.Delete(userID)
}
