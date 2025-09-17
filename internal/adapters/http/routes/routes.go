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
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/", homeHandler).Methods("GET")

	// Rutas de autenticación (públicas)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

	// Rutas protegidas (requieren JWT)
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)

	// Rutas de usuario
	users := protected.PathPrefix("/users").Subrouter()
	users.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")
	users.HandleFunc("/profile", userHandler.UpdateProfile).Methods("PUT")
	users.HandleFunc("", userHandler.ListUsers).Methods("GET")
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
	info := `
🔐 Auth Service API

Endpoints disponibles:

PÚBLICOS:
  GET  /health                     - Health check
  POST /api/v1/auth/register       - Registrar usuario
  POST /api/v1/auth/login          - Iniciar sesión
  POST /api/v1/auth/refresh        - Refrescar token

PROTEGIDOS (requieren Authorization: Bearer <token>):
  GET  /api/v1/users/profile       - Ver perfil
  PUT  /api/v1/users/profile       - Actualizar perfil
  GET  /api/v1/users               - Listar usuarios (paginado)
  GET  /api/v1/users/{id}          - Ver usuario por ID

USUARIOS DE PRUEBA:
  admin@example.com / admin123
  john@example.com  / john123
  jane@example.com  / jane123
  bob@example.com   / bob123
  alice@example.com / alice123

EJEMPLO DE USO:
  1. Login:
     curl -X POST http://localhost:8080/api/v1/auth/login \
       -H "Content-Type: application/json" \
       -d '{"email":"admin@example.com","password":"admin123"}'

  2. Usar token:
     curl -X GET http://localhost:8080/api/v1/users/profile \
       -H "Authorization: Bearer <tu-token-aqui>"
`

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(info))
}
