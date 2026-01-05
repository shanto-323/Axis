package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanto-323/axis/config"
	"github.com/shanto-323/axis/internal/database/mock"
	"github.com/shanto-323/axis/internal/database/postgres"
	"github.com/shanto-323/axis/internal/model"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/model/entity"
	"go.opentelemetry.io/otel/trace"
)

// It contains all methods that database should implement.
type Database interface {
	// Database specific methods
	Ping(ctx context.Context) error
	IsInitialized(ctx context.Context) bool
	Close() error

	// Other methods related to database operation
	CreateUser(ctx context.Context, user *dto.RegisterRequest) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	CreateConversationLog(ctx context.Context, cl *entity.ConversationLog) (*entity.ConversationLog, error)
	GetConversationLogHistory(ctx context.Context, userId uuid.UUID, queryDto *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error)
}

func New(cfg *config.Config, logger *zerolog.Logger, tracer trace.Tracer) (Database, error) {
	switch cfg.Database.Type {
	case "postgres":
		return postgres.New(cfg, logger, tracer)
	case "mock":
		return mock.New(cfg, logger)
	default:
		return nil, fmt.Errorf("no database found")
	}
}
