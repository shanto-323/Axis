package router

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/server"
	"github.com/shanto-323/axis/internal/server/handler"
	"github.com/shanto-323/axis/internal/server/middleware"
	v1 "github.com/shanto-323/axis/internal/server/router/v1"
)

const ApiVersion = "/api/v1"

func NewRouter(s *server.Server, h *handler.Handlers) *echo.Echo {
	middlewares := middleware.New(s)

	router := echo.New()

	router.HTTPErrorHandler = middlewares.GlobalErrorHandler

	router.Use(
		middlewares.CROS(),
		middlewares.RateLimitHit(),
		middlewares.Secure(),
		middlewares.TracerHandler(),
		middleware.RequestID(),
		middlewares.EnhanceContext(),
		middlewares.RequestLogger(),
		middlewares.Recover(),
	)

	registerSystemRoutes(router, h)

	r := router.Group(ApiVersion)
	v1.RegisterV1Routes(r, h, middlewares)
	return router
}
