package dto

import (
	"github.com/go-playground/validator"
	"github.com/shanto-323/axis/internal/model"
)

type ConversationHistoryQuery struct {
	Page  *int `query:"page" validate:"omitempty,min=1"`
	Limit *int `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (l *ConversationHistoryQuery) Validate() error {
	if err := validator.New().Struct(l); err != nil {
		return err
	}

	if l.Page == nil {
		defaultPage := 1
		l.Page = &defaultPage
	}
	if l.Limit == nil {
		defaultLimit := 10
		l.Limit = &defaultLimit
	}
	return nil
}

type ConversationLogResponse struct {
	model.BaseLV

	TextQuery    string   `json:"query"`
	ResponseText string   `json:"response_text"`
	TimeTaken    int      `json:"time_taken"`
}
