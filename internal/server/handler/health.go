package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/server"
	"github.com/shanto-323/axis/internal/server/middleware"
)

const (
	Healthy   = "healthy"
	Unhealthy = "unhealthy"
)

type HealthHandler struct {
	server *server.Server
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		server: s,
	}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	start := time.Now()
	logger := middleware.GetLogger(c).With().
		Str("operation", "health_check").
		Logger()

	isHealthy := true

	if !h.server.Database.IsInitialized(context.Background()) {
		return fmt.Errorf("database not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbStart := time.Now()
	if err := h.server.Database.Ping(ctx); err != nil {
		isHealthy = false
		logger.Error().
			Str("check_type", "postgres").
			Str("operation", "health_check").
			Str("error_type", "postgres_unhealthy").
			Int64("response_time_ms", time.Since(dbStart).Milliseconds()).
			Str("error_message", err.Error()).
			Msg("HealthCheckError")
	} else {
		logger.Info().
			Dur("response_time", time.Since(dbStart)).
			Msg("database health check passed")
	}

	if !isHealthy {
		logger.Error().
			Str("check_type", "overall").
			Str("operation", "health_check").
			Str("error_type", "overall_unhealthy").
			Int64("total_duration_ms", time.Since(start).Milliseconds()).
			Msg("HealthCheckError")
		return c.JSON(http.StatusServiceUnavailable, nil)
	}

	logger.Info().
		Dur("total_duration", time.Since(start)).
		Msg("health check passed")

	err := c.JSON(http.StatusOK, nil)
	if err != nil {
		logger.Error().
			Str("check_type", "response").
			Str("operation", "health_check").
			Str("error_type", "json_response_error").
			Str("error_message", err.Error()).
			Msg("HealthCheckError")
		return fmt.Errorf("failed to write JSON response: %w", err)
	}

	return nil
}
