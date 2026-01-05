package service

import (
	"github.com/shanto-323/axis/internal/server"
)

type Services struct {
	Auth AuthService
	Chat ChatService
}

func New(s *server.Server) *Services {
	return &Services{
		Auth: NewAuthService(s.Config, s.Database, s.Tracer.Tracer),
		Chat: NewChatService(s.LLM, s.Database, s.Tracer.Tracer),
	}
}
