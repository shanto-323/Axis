package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/server"
	"github.com/shanto-323/axis/internal/server/middleware"
	"github.com/shanto-323/axis/internal/validation"
	"go.opentelemetry.io/otel/attribute"
)

type Handler struct {
	server *server.Server
}

func NewHandler(s *server.Server) *Handler {
	return &Handler{
		server: s,
	}
}

type HandlerFunc[Req validation.Validatable, Res any] func(c echo.Context, req Req) (Res, error)

type HandleNoResponseFunc[Req validation.Validatable] func(c echo.Context, req Req) error

type ResponseHandler interface {
	Handle(c echo.Context, result any) error
	GetOperation() string
}

type JSONResponseHandler struct {
	status int
}

func (h JSONResponseHandler) Handle(c echo.Context, result any) error {
	return c.JSON(h.status, result)
}

func (h JSONResponseHandler) GetOperation() string {
	return "handler"
}

type NoResponseHandler struct {
	status int
}

func (h NoResponseHandler) Handle(c echo.Context, result any) error {
	return c.NoContent(h.status)
}

func (h NoResponseHandler) GetOperation() string {
	return "handler_no_response"
}

func handleRequest[Req validation.Validatable](
	h *Handler,
	c echo.Context,
	req Req,
	handler func(c echo.Context, req Req) (any, error),
	responseHandler ResponseHandler,
) error {
	start := time.Now()
	method := c.Request().Method
	path := c.Path()

	ctx := c.Request().Context()

	ctx, span := h.server.Tracer.Tracer.Start(ctx, "handler")
	defer span.End()

	c.SetRequest(c.Request().WithContext(ctx))

	// Get context-enhanced logger
	loggerBuilder := middleware.GetLogger(c).With().
		Str("operation", responseHandler.GetOperation()).
		Str("method", method).
		Str("path", path)

	logger := loggerBuilder.Logger()

	validationStart := time.Now()
	if err := validation.BindAndValidate(c, req); err != nil {
		validationDuration := time.Since(validationStart)

		logger.Error().
			Err(err).
			Dur("validation_duration", validationDuration).
			Msg("request validation failed")
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("validation.status", "failed"),
			attribute.Int64("validation.duration_ms", validationDuration.Milliseconds()),
		)

		return err
	}

	validationDuration := time.Since(validationStart)

	logger.Debug().
		Dur("validation_duration", validationDuration).
		Msg("request validation successful")
	span.SetAttributes(
		attribute.String("validation.status", "success"),
		attribute.Int64("validation.duration_ms", validationDuration.Milliseconds()),
	)

	handlerStart := time.Now()
	result, err := handler(c, req)
	handlerDuration := time.Since(handlerStart)

	if err != nil {
		totalDuration := time.Since(start)

		logger.Error().
			Err(err).
			Dur("handler_duration", handlerDuration).
			Dur("total_duration", totalDuration).
			Msg("handler execution failed")

		span.RecordError(err)
		span.SetAttributes(
			attribute.String("handler.status", "failed"),
			attribute.Int64("handler.duration_ms", handlerDuration.Milliseconds()),
			attribute.Int64("total.duration_ms", validationDuration.Milliseconds()),
		)

		return err
	}

	totalDuration := time.Since(start)

	logger.Info().
		Dur("handler_duration", handlerDuration).
		Dur("validation_duration", validationDuration).
		Dur("total_duration", totalDuration).
		Msg("request completed successfully")

	span.SetAttributes(
		attribute.String("handler.status", "success"),
		attribute.Int64("handler.duration_ms", handlerDuration.Milliseconds()),
		attribute.Int64("total.duration_ms", totalDuration.Milliseconds()),
	)

	return responseHandler.Handle(c, result)
}

func Handle[Req validation.Validatable, Res any](
	h *Handler,
	handler HandlerFunc[Req, Res],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(h, c, req, func(c echo.Context, req Req) (any, error) {
			return handler(c, req)
		}, JSONResponseHandler{status: status})
	}
}

func HandleNoResponse[Req validation.Validatable](
	h *Handler,
	handler HandleNoResponseFunc[Req],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(h, c, req, func(c echo.Context, req Req) (any, error) {
			return nil, handler(c, req)
		}, NoResponseHandler{status: status})
	}
}
