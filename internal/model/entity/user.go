package entity

import "github.com/shanto-323/axis/internal/model"

type User struct {
	model.Base

	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}
