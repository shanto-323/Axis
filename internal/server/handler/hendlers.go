package handler

import (
	"github.com/shanto-323/axis/internal/server"
	"github.com/shanto-323/axis/internal/service"
)

type Handlers struct {
	Auth    *AuthHandler
	Chat    *ChatHandler
	OpenAPI *OpenAPIHandler
	Health  *HealthHandler
}

func New(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Auth:    NewAuthHandler(s, services.Auth),
		Chat:    NewChatHandler(s, services.Chat),
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(),
	}
}
