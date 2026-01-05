package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/model"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/model/entity"
	"github.com/shanto-323/axis/internal/server"
	"github.com/shanto-323/axis/internal/service"
)

type ChatHandler struct {
	*Handler
	service service.ChatService
}

func NewChatHandler(s *server.Server, service service.ChatService) *ChatHandler {
	return &ChatHandler{
		Handler: NewHandler(s),
		service: service,
	}
}

func (h *ChatHandler) ModelHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := h.service.AvailableModels(c)
		return c.JSON(http.StatusOK, resp)
	}
}

func (h *ChatHandler) ChatHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return Handle(
			h.Handler,
			func(c echo.Context, req *dto.ChatRequest) (*entity.ConversationLog, error) {
				return h.service.Chat(c, req)
			},
			http.StatusOK,
			&dto.ChatRequest{},
		)(c)
	}
}

func (h *ChatHandler) ChatHistoryHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return Handle(
			h.Handler,
			func(c echo.Context, req *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error) {
				return h.service.ChatHistory(c, req)
			},
			http.StatusOK,
			&dto.ConversationHistoryQuery{},
		)(c)
	}
}
