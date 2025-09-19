package service

import (
	"encoding/json"
	"net/http"

	"github.com/danielgalindoj/auth-service/internal/core/domain/model"
	"github.com/danielgalindoj/auth-service/internal/core/services"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Validaciones básicas
	if req.Email == "" {
		http.Error(w, `{"error":"Email requerido"}`, http.StatusBadRequest)
		return
	}
	if req.Password == "" || len(req.Password) < 6 {
		http.Error(w, `{"error":"Contraseña debe tener al menos 6 caracteres"}`, http.StatusBadRequest)
		return
	}

	firstName := ""
	if req.FirstName != nil {
		firstName = *req.FirstName
	}

	lastName := ""
	if req.LastName != nil {
		lastName = *req.LastName
	}

	username := ""
	if req.Username != nil {
		username = *req.Username
	}

	user, err := h.authService.Register(req.Email, req.Password, firstName, lastName, username)
	if err != nil {
		logrus.Error("Error en registro: ", err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.UserResponse{
		User:    user,
		Message: "Usuario registrado exitosamente",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Validaciones básicas
	if req.Email == "" {
		http.Error(w, `{"error":"Email requerido"}`, http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		http.Error(w, `{"error":"Contraseña requerida"}`, http.StatusBadRequest)
		return
	}

	loginResponse, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		logrus.Warn("Error en login: ", err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, `{"error":"Refresh token requerido"}`, http.StatusBadRequest)
		return
	}

	tokenPair, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		logrus.Warn("Error refrescando token: ", err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":         tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_at":    tokenPair.ExpiresAt,
	})
}
