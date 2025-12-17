package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"fullName"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
}

type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

type AchievementHistoryEntry struct {
	ID            uuid.UUID  `json:"id"`
	AchievementID uuid.UUID  `json:"achievement_id"`
	Status        string     `json:"status"`
	ChangedBy     *uuid.UUID `json:"changed_by"`
	ChangedByName *string    `json:"changed_by_name"`
	RejectionNote *string    `json:"rejection_note"`
	CreatedAt     time.Time  `json:"created_at"`
}
