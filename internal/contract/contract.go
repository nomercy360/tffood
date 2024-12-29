package contract

import "time"

type Error struct {
	Message string `json:"message"`
}

type UserAuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID                   int64     `json:"id"`
	FirstName            *string   `json:"first_name"`
	LastName             *string   `json:"last_name"`
	Username             string    `json:"username"`
	ChatID               int64     `json:"chat_id"`
	LanguageCode         *string   `json:"language"`
	IsPremium            bool      `json:"is_premium"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	LastSeenAt           time.Time `json:"last_seen_at"`
	NotificationsEnabled bool      `json:"notifications_enabled"`
	AvatarURL            *string   `json:"avatar_url"`
	Title                *string   `json:"title"`
}
