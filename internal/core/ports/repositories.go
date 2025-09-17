package ports

import (
	"github.com/danielgalindoj/auth-service/internal/core/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*domain.User, error)
	Count() (int64, error)
}
