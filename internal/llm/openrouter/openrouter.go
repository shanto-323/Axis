package openrouter

import (
	"context"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/rs/zerolog"
	"github.com/shanto-323/axis/config"
	"github.com/shanto-323/axis/internal/errs"
	"github.com/shanto-323/axis/internal/llm"
	"github.com/shanto-323/axis/internal/model/dto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Openrouter struct {
	logger *zerolog.Logger
	config *config.Config
	tracer trace.Tracer

	llmModels llm.LLMModels
}

func NewOpenrouter(cfg *config.Config, log *zerolog.Logger, tracer trace.Tracer) *Openrouter {
	// All available model
	models := map[string]string{
		"openai/gpt-120b": "openai/gpt-oss-120b:free",
		"llama-70b":       "meta-llama/llama-3.3-70b-instruct:free",
		"nemotron-30b":    "nvidia/nemotron-3-nano-30b-a3b:free",
		"nemotron-12b":    "nvidia/nemotron-nano-12b-v2-vl:free",
		"qwen3":           "qwen/qwen3-coder:free",
		"allenai-32b":     "allenai/olmo-3.1-32b-think:free",
		"xiaomi-flash":    "xiaomi/mimo-v2-flash:free",
		"mistralai":       "mistralai/devstral-2512:free",
		"deepseek-nex":    "nex-agi/deepseek-v3.1-nex-n1:free",
		"tngtech":         "tngtech/tng-r1t-chimera:free",
		"kat-coder":       "kwaipilot/kat-coder-pro:free",
	}

	return &Openrouter{
		logger:    log,
		config:    cfg,
		tracer:    tracer,
		llmModels: models,
	}
}

func (o *Openrouter) AvailableModels(ctx context.Context) *[]dto.LLMModel {
	_, span := o.tracer.Start(ctx, "event.models")
	defer span.End()

	models := []dto.LLMModel{}
	for k, v := range o.llmModels {
		model := dto.LLMModel{
			Name:  k,
			Model: v,
		}
		models = append(models, model)
	}

	span.SetAttributes(
		attribute.String("get all models", "success"),
	)

	return &models
}

func (o *Openrouter) GenerateResponse(ctx context.Context, request *dto.ChatRequest) (*dto.ConversationLogResponse, error) {
	ctx, span := o.tracer.Start(ctx, "event.llm_response")
	defer span.End()

	startTime := time.Now()

	model, ok := o.llmModels[request.Model]
	if !ok {
		code := "INVALID_MODEL_NAME"
		return nil, errs.NewNotFoundError("no such model found :"+request.Model, true, &code)
	}

	responseString, _, err := o.response(
		ctx,
		request.Message,
		model,
	)
	if err != nil {
		o.logger.Error().
			Err(err).
			Str("event", "llm-response").
			Msg("response gen failed")

		span.RecordError(err)
		return nil, errs.NewInternalServerError()
	}

	totalTime := int(time.Since(startTime).Seconds())

	o.logger.Info().
		Str("event", "llm-response").
		Int("time", totalTime).
		Msg("success")

	response := dto.ConversationLogResponse{
		TextQuery:    request.Message,
		ResponseText: responseString,
		TimeTaken:    totalTime,
	}

	response.LLMModelName = request.Model
	response.Timestamp = time.Now()

	return &response, nil
}

func (o *Openrouter) response(ctx context.Context, query string, model string) (string, int, error) {
	client := openai.NewClient(
		option.WithBaseURL(o.config.AiManage.Provider),
		option.WithAPIKey(o.config.AiManage.ApiKey),
	)

	start := time.Now()
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(query),
		},
		Model: model,
	})
	if err != nil {
		return "", 0, err
	}

	exicutionTime := int(time.Since(start).Seconds())

	o.logger.Info().
		Str("event", "llm-response").
		Int("time", exicutionTime).
		Msg("successful")

	return resp.Choices[0].Message.Content, exicutionTime, nil
}
