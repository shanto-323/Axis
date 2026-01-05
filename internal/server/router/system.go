package router

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/server/handler"
)

func registerSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.Static("/static", "static")
	r.GET("/docs", h.OpenAPI.ServeOpenAPIUI)

	r.GET("/health", h.Health.CheckHealth)

}
