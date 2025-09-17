package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danielgalindoj/auth-service/internal/adapters/http/handlers"
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

	// Rutas públicas
	api := r.PathPrefix("/api/" + config.APIVersion).Subrouter()

	// Ruta de salud
	// @Summary Verifica estado del servicio
	// @Description Retorna estado, versión y timestamp del servicio.
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// @Summary Página principal
	// @Description Información general del Auth Service.
	// @Tags Health
	// @Produce plain
	// @Success 200 {string} string "Información del servicio"
	// @Router / [get]
	r.HandleFunc("/", homeHandler).Methods("GET")

	// Rutas de autenticación (públicas)
	auth := api.PathPrefix("/auth").Subrouter()

	// @Summary Registrar usuario
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param data body map[string]string true "Email y Password"
	// @Success 201 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Router /auth/register [post]
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")

	// @Summary Login de usuario
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param data body map[string]string true "Email y Password"
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]string
	// @Router /auth/login [post]
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")

	// @Summary Refrescar token
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param data body map[string]string true "Refresh token"
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]string
	// @Router /auth/refresh [post]
	auth.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

	// Rutas protegidas
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)

	// Rutas de usuario
	users := protected.PathPrefix("/users").Subrouter()

	// @Summary Ver perfil del usuario autenticado
	// @Tags Users
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]string
	// @Router /users/profile [get]
	users.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")

	// @Summary Actualizar perfil del usuario autenticado
	// @Tags Users
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param data body map[string]string true "Datos de perfil"
	// @Success 200 {object} map[string]string
	// @Failure 401 {object} map[string]string
	// @Router /users/profile [put]
	users.HandleFunc("/profile", userHandler.UpdateProfile).Methods("PUT")

	// @Summary Listar usuarios
	// @Tags Users
	// @Security BearerAuth
	// @Produce json
	// @Param page query int false "Número de página"
	// @Param limit query int false "Resultados por página"
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]string
	// @Router /users [get]
	users.HandleFunc("", userHandler.ListUsers).Methods("GET")

	// @Summary Obtener usuario por ID
	// @Tags Users
	// @Security BearerAuth
	// @Produce json
	// @Param id path string true "ID del usuario"
	// @Success 200 {object} map[string]interface{}
	// @Failure 404 {object} map[string]string
	// @Failure 401 {object} map[string]string
	// @Router /users/{id} [get]
	users.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"service":   "auth-service",
		"version":   "1.0.0",
		"timestamp": "2024-01-01T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	info := `Auth Service API`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(info))
}
