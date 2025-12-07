package model

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	ISActive     bool      `json:"is_active"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"update_at"`
}