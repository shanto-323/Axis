package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/shanto-323/axis/internal/model"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/model/entity"
)

func (db *DB) CreateConversationLog(ctx context.Context, cl *entity.ConversationLog) (*entity.ConversationLog, error) {
	cl.ID = uuid.New()

	idString := cl.ID.String()

	db.pool[idString] = &cl
	db.logger.Info().
		Str("event", "new_log").
		Str("user_id", cl.UserID.String()).
		Str("llm_model", cl.LLMModelName).
		Msg("new user created")

	return cl, nil
}

func (db *DB) GetConversationLogHistory(ctx context.Context, userId uuid.UUID, query *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error) {
	return nil, nil
}
