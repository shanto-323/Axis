package dto

import "github.com/go-playground/validator"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *LoginRequest) Validate() error {
	return validator.New().Struct(r)
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterRequest) Validate() error {
	return validator.New().Struct(r)
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}
