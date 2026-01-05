package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/internal/database"
	"github.com/shanto-323/axis/internal/errs"
	"github.com/shanto-323/axis/internal/llm"
	"github.com/shanto-323/axis/internal/model"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/model/entity"
	"go.opentelemetry.io/otel/trace"
)

type ChatService interface {
	AvailableModels(c echo.Context) *[]dto.LLMModel
	Chat(c echo.Context, payload *dto.ChatRequest) (*entity.ConversationLog, error)
	ChatHistory(c echo.Context, payload *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error)
}

type chatService struct {
	db     database.Database
	llm    llm.LLM
	tracer trace.Tracer
}

func NewChatService(llm llm.LLM, db database.Database, tracer trace.Tracer) *chatService {
	return &chatService{
		db:     db,
		llm:    llm,
		tracer: tracer,
	}
}

func (s *chatService) AvailableModels(c echo.Context) *[]dto.LLMModel {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 60*time.Second)
	defer cancel()

	return s.llm.AvailableModels(ctx)
}

func (s *chatService) Chat(c echo.Context, payload *dto.ChatRequest) (*entity.ConversationLog, error) {
	ctx := c.Request().Context()

	ctx, span := s.tracer.Start(ctx, "service")
	defer span.End()

	c.SetRequest(c.Request().WithContext(ctx))

	userId, ok := c.Get("id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	llmResponse, err := s.llm.GenerateResponse(ctx, payload)
	if err != nil {
		return nil, err
	}

	cLog := entity.ConversationLog{
		BaseLV:       llmResponse.BaseLV,
		UserID:       userId,
		TextQuery:    llmResponse.TextQuery,
		ResponseText: llmResponse.ResponseText,
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	return s.db.CreateConversationLog(ctx, &cLog)
}

func (s *chatService) ChatHistory(c echo.Context, payload *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error) {
	ctx := c.Request().Context()

	ctx, span := s.tracer.Start(ctx, "service")
	defer span.End()

	c.SetRequest(c.Request().WithContext(ctx))

	userId, ok := c.Get("id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	return s.db.GetConversationLogHistory(ctx, userId, payload)
}
