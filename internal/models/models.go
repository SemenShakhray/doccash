package models

import "github.com/google/uuid"

type User struct {
	UserID       uuid.UUID
	Login        string
	PasswordHash []byte
	CreatedAt    int64
}
