package model

import "time"

type BaseLLMModel struct {
	LLMModelName string `json:"llm_model_name" db:"llm_model_name"`
}

type BaseTimestamp struct {
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type BaseLV struct {
	BaseLLMModel
	BaseTimestamp
}
