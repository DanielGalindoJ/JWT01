package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/danielgalindoj/auth-service/internal/core/domain/model"
	"github.com/danielgalindoj/auth-service/internal/core/ports"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (id, email, username, password_hash, first_name, last_name, 
			email_verified, is_active, provider, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.EmailVerified,
		user.IsActive,
		user.Provider,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		logrus.Error("Error creando usuario en base de datos: ", err)
		return fmt.Errorf("error creando usuario: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
			avatar_url, email_verified, is_active, provider, provider_id, 
			last_login_at, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = true
	`

	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.EmailVerified,
		&user.IsActive,
		&user.Provider,
		&user.ProviderID,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
			avatar_url, email_verified, is_active, provider, provider_id, 
			last_login_at, created_at, updated_at
		FROM users WHERE email = $1 AND is_active = true
	`

	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.EmailVerified,
		&user.IsActive,
		&user.Provider,
		&user.ProviderID,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario por email: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
			avatar_url, email_verified, is_active, provider, provider_id, 
			last_login_at, created_at, updated_at
		FROM users WHERE username = $1 AND is_active = true
	`

	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.EmailVerified,
		&user.IsActive,
		&user.Provider,
		&user.ProviderID,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo usuario por username: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(user *model.User) error {
	query := `
		UPDATE users SET 
			email = $2, username = $3, password_hash = $4, first_name = $5, 
			last_name = $6, avatar_url = $7, email_verified = $8, 
			last_login_at = $9, updated_at = $10
		WHERE id = $1
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.AvatarURL,
		user.EmailVerified,
		user.LastLoginAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error actualizando usuario: %w", err)
	}

	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = $2 WHERE id = $1`

	_, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("error eliminando usuario: %w", err)
	}

	return nil
}

func (r *userRepository) List(limit, offset int) ([]*model.User, error) {
	query := `
		SELECT id, email, username, first_name, last_name, avatar_url, 
			email_verified, is_active, provider, last_login_at, created_at, updated_at
		FROM users 
		WHERE is_active = true 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo lista de usuarios: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.AvatarURL,
			&user.EmailVerified,
			&user.IsActive,
			&user.Provider,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando usuario: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_active = true`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error contando usuarios: %w", err)
	}

	return count, nil
}
