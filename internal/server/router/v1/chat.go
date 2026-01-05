package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/server/handler"
	"github.com/shanto-323/axis/internal/server/middleware"
)

func registerChatRoute(r *echo.Group, h *handler.Handlers, m *middleware.Middlewares) {
	chatRoute := r.Group("/chat")
	{
		chatRoute.Use(m.RequireAuth())
		chatRoute.POST("", h.Chat.ChatHandler())
		chatRoute.GET("/models", h.Chat.ModelHandler())
		chatRoute.POST("/history", h.Chat.ChatHistoryHandler())
	}
}
