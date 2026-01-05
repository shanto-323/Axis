package dto

import (
	"github.com/go-playground/validator"
)

type ChatRequest struct {
	Model   string `json:"model"`
	Message string `json:"message" validate:"required"`
}

func (r *ChatRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return err
	}

	if r.Model == "" {
		r.Model = "llama_70b"
	}

	return nil
}
