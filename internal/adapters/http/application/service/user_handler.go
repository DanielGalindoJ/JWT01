package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/danielgalindoj/auth-service/internal/adapters/http/middleware"
	"github.com/danielgalindoj/auth-service/internal/core/domain/model"
	"github.com/danielgalindoj/auth-service/internal/core/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, `{"error":"Usuario no encontrado en contexto"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.UserResponse{
		User:    user,
		Message: "Perfil del usuario",
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, `{"error":"ID de usuario requerido"}`, http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, `{"error":"ID de usuario inválido"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		http.Error(w, `{"error":"Usuario no encontrado"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parámetros de paginación
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	users, total, err := h.userService.List(page, limit)
	if err != nil {
		logrus.Error("Error obteniendo usuarios: ", err)
		http.Error(w, `{"error":"Error obteniendo usuarios"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"users":       users,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, `{"error":"Usuario no encontrado en contexto"}`, http.StatusInternalServerError)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userService.Update(user.ID, updates)
	if err != nil {
		logrus.Error("Error actualizando perfil: ", err)
		http.Error(w, `{"error":"Error actualizando perfil"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.UserResponse{
		User:    updatedUser,
		Message: "Perfil actualizado exitosamente",
	})
}
