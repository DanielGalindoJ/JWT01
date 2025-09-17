package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgalindoj/auth-service/internal/core/domain"
	"github.com/danielgalindoj/auth-service/internal/core/services"
	"github.com/sirupsen/logrus"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener token del header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Token requerido"}`, http.StatusUnauthorized)
				return
			}

			// Verificar formato Bearer
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, `{"error":"Formato de token inválido"}`, http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]

			// Validar token
			user, err := authService.ValidateToken(tokenString)
			if err != nil {
				logrus.Warn("Token inválido: ", err)
				http.Error(w, `{"error":"Token inválido"}`, http.StatusUnauthorized)
				return
			}

			// Agregar usuario al contexto
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext extrae el usuario del contexto
func GetUserFromContext(r *http.Request) (*domain.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*domain.User)
	return user, ok
}
