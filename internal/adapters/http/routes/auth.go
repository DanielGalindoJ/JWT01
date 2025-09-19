package routes

import (
	handlers "github.com/danielgalindoj/auth-service/internal/adapters/http/application/service"
	"github.com/gorilla/mux"
)

// setupAuthRoutes configura las rutas de autenticación
func setupAuthRoutes(api *mux.Router, authHandler *handlers.AuthHandler) {
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
}
