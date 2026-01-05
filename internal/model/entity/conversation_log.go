package entity

import (
	"github.com/google/uuid"
	"github.com/shanto-323/axis/internal/model"
)

type ConversationLog struct {
	model.BaseId
	model.BaseLV

	UserID       uuid.UUID `db:"user_id" json:"user_id"`
	TextQuery    string    `db:"text_query" json:"query"`
	ResponseText string    `db:"response_text" json:"response_text"`
}
