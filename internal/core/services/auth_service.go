package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/danielgalindoj/auth-service/internal/config"
	"github.com/danielgalindoj/auth-service/internal/core/domain"
	"github.com/danielgalindoj/auth-service/internal/core/ports"
	"github.com/danielgalindoj/auth-service/pkg/hash"
	"github.com/danielgalindoj/auth-service/pkg/jwt"
)

type AuthService struct {
	userRepo ports.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo ports.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *AuthService) Register(email, password, firstName, lastName, username string) (*domain.User, error) {
	// Verificar si el usuario ya existe
	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, fmt.Errorf("usuario con email %s ya existe", email)
	}

	// Verificar username si se proporciona
	if username != "" {
		if _, err := s.userRepo.GetByUsername(username); err == nil {
			return nil, fmt.Errorf("usuario con username %s ya existe", username)
		}
	}

	// Hashear contraseña
	hashedPassword, err := hash.HashPassword(password, s.config.BCryptCost)
	if err != nil {
		return nil, fmt.Errorf("error hasheando contraseña: %w", err)
	}

	// Crear usuario
	user := &domain.User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  &hashedPassword,
		EmailVerified: false,
		IsActive:      true,
		Provider:      "local",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if firstName != "" {
		user.FirstName = &firstName
	}
	if lastName != "" {
		user.LastName = &lastName
	}
	if username != "" {
		user.Username = &username
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	logrus.WithField("email", email).Info("Usuario registrado exitosamente")

	// Limpiar password hash antes de retornar
	user.PasswordHash = nil
	return user, nil
}

func (s *AuthService) Login(email, password string) (*domain.LoginResponse, error) {
	// Buscar usuario por email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Verificar contraseña
	if user.PasswordHash == nil || !hash.CheckPasswordHash(password, *user.PasswordHash) {
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, fmt.Errorf("cuenta desactivada")
	}

	// Generar tokens
	username := ""
	if user.Username != nil {
		username = *user.Username
	}

	tokenPair, err := jwt.GenerateTokenPair(
		user.ID,
		user.Email,
		username,
		s.config.JWTSecret,
		s.config.JWTExpiration,
		s.config.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, fmt.Errorf("error generando tokens: %w", err)
	}

	// Actualizar last login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user); err != nil {
		logrus.Warn("Error actualizando last login: ", err)
	}

	logrus.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("Login exitoso")

	// Limpiar password hash antes de retornar
	user.PasswordHash = nil

	return &domain.LoginResponse{
		User:         user,
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

func (s *AuthService) RefreshToken(refreshTokenString string) (*jwt.TokenPair, error) {
	// Validar refresh token
	claims, err := jwt.ValidateToken(refreshTokenString, s.config.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("refresh token inválido: %w", err)
	}

	// Buscar usuario
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("cuenta desactivada")
	}

	// Generar nuevos tokens
	username := ""
	if user.Username != nil {
		username = *user.Username
	}

	tokenPair, err := jwt.GenerateTokenPair(
		user.ID,
		user.Email,
		username,
		s.config.JWTSecret,
		s.config.JWTExpiration,
		s.config.JWTRefreshExpiration,
	)
	if err != nil {
		return nil, fmt.Errorf("error generando nuevos tokens: %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*domain.User, error) {
	claims, err := jwt.ValidateToken(tokenString, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("cuenta desactivada")
	}

	// Limpiar password hash
	user.PasswordHash = nil
	return user, nil
}
