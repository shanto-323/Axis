package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/errs"
	"github.com/shanto-323/axis/internal/server"
	"github.com/shanto-323/axis/pkg"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{server: s}
}

func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("access_token")
			if err != nil {
				return errs.NewUnauthorizedError("missing access token", false)
			}

			claims, err := pkg.ValidateToken(m.server.Config, cookie.Value)
			if err != nil {
				return errs.NewForbiddenError("invalid or expired access token", false)
			}

			c.Set("id", claims.ID)
			ctx := c.Request().Context()

			newCtxWithID := context.WithValue(ctx, "id", claims.ID)
			c.SetRequest(c.Request().WithContext(newCtxWithID))

			return next(c)
		}
	}
}
