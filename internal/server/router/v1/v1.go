package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/server/handler"
	"github.com/shanto-323/axis/internal/server/middleware"
)

func RegisterV1Routes(r *echo.Group, h *handler.Handlers, m *middleware.Middlewares) {
	registerAuthRoutes(r, h)

	registerChatRoute(r, h, m)
}
