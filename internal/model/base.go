package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseId struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type BaseCreatedAt struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BaseUpdatedAt struct {
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Base struct {
	BaseId
	BaseCreatedAt
	BaseUpdatedAt
}

