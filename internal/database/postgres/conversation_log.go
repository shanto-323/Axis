package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shanto-323/axis/internal/errs"
	"github.com/shanto-323/axis/internal/model"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/model/entity"
)

func (db *DB) CreateConversationLog(ctx context.Context, cl *entity.ConversationLog) (*entity.ConversationLog, error) {
	query := `
		INSERT INTO conversation_logs (
			user_id,
			text_query,
			response_text,
			llm_model_name
		)
		VALUES (
			@user_id,
			@text_query,
			@response_text,
			@llm_model_name
		)	
		RETURNING 
			* 
	`

	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id":        cl.UserID,
		"text_query":     cl.TextQuery,
		"response_text":  cl.ResponseText,
		"llm_model_name": cl.LLMModelName,
	}).Scan(
		&cl.ID,
		&cl.UserID,
		&cl.TextQuery,
		&cl.ResponseText,
		&cl.LLMModelName,
		&cl.Timestamp,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewInternalServerError()
		}
		return nil, err
	}

	return cl, nil
}

func (db *DB) GetConversationLogHistory(
	ctx context.Context,
	userId uuid.UUID,
	queryDto *dto.ConversationHistoryQuery,
) (*model.PaginatedResponse[entity.ConversationLog], error) {

	x := *queryDto.Page
	y := *queryDto.Limit
	offset := (x - 1) * y

	query := `
		SELECT 
			*
		FROM 
			conversation_logs
		WHERE
			user_id=@user_id
		ORDER BY
			timestamp DESC
		LIMIT @limit
		OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userId,
		"limit":   y,
		"offset":  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query")
	}

	logs, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.ConversationLog])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[entity.ConversationLog]{
				Data:       []entity.ConversationLog{},
				Page:       *queryDto.Page,
				Limit:      *queryDto.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows")
	}

	// total count
	count := `
		SELECT
			COUNT(*)
		FROM
			conversation_logs	
		WHERE
			user_id=@user_id
	`

	countArgs := pgx.NamedArgs{
		"user_id": userId,
	}

	var total int
	err = db.pool.QueryRow(ctx, count, countArgs).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count")
	}

	return &model.PaginatedResponse[entity.ConversationLog]{
		Data:       logs,
		Page:       *queryDto.Page,
		Limit:      *queryDto.Limit,
		Total:      total,
		TotalPages: (total + *queryDto.Limit - 1) / *queryDto.Limit,
	}, nil
}
