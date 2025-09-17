package http

import (
	"github.com/danielgalindoj/auth-service/internal/config"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}
