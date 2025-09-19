package ports

import (
	"github.com/danielgalindoj/auth-service/internal/core/domain/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	//Persistence
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uuid.UUID) error

	//Query
	GetByID(id uuid.UUID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	List(limit, offset int) ([]*model.User, error)
	Count() (int64, error)
}
