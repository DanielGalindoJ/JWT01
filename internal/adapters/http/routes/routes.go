package routes

import (
	handlers "github.com/danielgalindoj/auth-service/internal/adapters/http/application/service"
	"github.com/danielgalindoj/auth-service/internal/adapters/http/middleware"
	"github.com/danielgalindoj/auth-service/internal/config"
	"github.com/danielgalindoj/auth-service/internal/core/services"
	"github.com/gorilla/mux"
)

// SetupRoutes configura todas las rutas de la API.
// @title Auth Service API
// @version 1.0
// @description 🔐 Servicio de autenticación y gestión de usuarios.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func SetupRoutes(authService *services.AuthService, userService *services.UserService, config *config.Config) *mux.Router {
	r := mux.NewRouter()

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Middleware de autenticación
	authMiddleware := middleware.AuthMiddleware(authService)

	// Rutas base
	api := r.PathPrefix("/api/" + config.APIVersion).Subrouter()

	// Configurar rutas de health
	setupHealthRoutes(r)

	// Configurar rutas de autenticación
	setupAuthRoutes(api, authHandler)

	// Configurar rutas protegidas de usuarios
	setupUserRoutes(api, authMiddleware, userHandler)

	return r
}
