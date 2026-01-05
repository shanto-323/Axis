package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shanto-323/axis/config"
	"github.com/shanto-323/axis/internal/model/entity"
)

type MockDb map[string]any

type DB struct {
	pool   MockDb
	logger *zerolog.Logger
}

func New(config *config.Config, logger *zerolog.Logger) (*DB, error) {
	logger.Debug().Msg("mock database started")

	return &DB{
		pool:   make(MockDb),
		logger: logger,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error         { return nil }
func (db *DB) IsInitialized(ctx context.Context) bool { return true }
func (db *DB) Close() error                           { return nil }
func (db *DB) GetHistoryForLLM(ctx context.Context, userId uuid.UUID) (*[]entity.ConversationLog, error) {
	return nil, nil
}
