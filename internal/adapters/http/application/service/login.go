package service

import "github.com/danielgalindoj/auth-service/internal/core/services"

type LoginHandler struct {
	loginService *services.AuthService
}

func NewLoginHandler(loginService *services.AuthService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}
