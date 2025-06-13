package model

import "time"

type User struct {
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	PasswordHash string    `json:"-"`
	Nickname    string    `json:"nickname"`
	AvatarURL   string    `json:"avatar_url"`
	ContactInfo string    `json:"contact_info"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}