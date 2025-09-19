package routes

import (
	handlers "github.com/danielgalindoj/auth-service/internal/adapters/http/application/service"
	"github.com/gorilla/mux"
)

// setupUserRoutes configura las rutas protegidas de usuarios
func setupUserRoutes(api *mux.Router, authMiddleware mux.MiddlewareFunc, userHandler *handlers.UserHandler) {
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
}
