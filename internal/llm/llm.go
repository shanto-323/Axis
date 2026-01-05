package llm

import (
	"context"

	"github.com/shanto-323/axis/internal/model/dto"
)

type LLMModels map[string]string

type LLM interface {
	GenerateResponse(ctx context.Context, request *dto.ChatRequest) (*dto.ConversationLogResponse, error)
	AvailableModels(ctx context.Context) *[]dto.LLMModel
}
