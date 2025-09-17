package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	Username      *string    `json:"username,omitempty"`
	PasswordHash  *string    `json:"-"` // No incluir en JSON
	FirstName     *string    `json:"first_name,omitempty"`
	LastName      *string    `json:"last_name,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	EmailVerified bool       `json:"email_verified"`
	IsActive      bool       `json:"is_active"`
	Provider      string     `json:"provider"`
	ProviderID    *string    `json:"provider_id,omitempty"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateUserRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Username  *string `json:"username,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User         *User  `json:"user"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserResponse struct {
	*User
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
