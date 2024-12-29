package db

import (
	"context"
	"time"
)

type User struct {
	ID                   int64     `db:"id"`
	FirstName            *string   `db:"first_name"`
	LastName             *string   `db:"last_name"`
	Username             string    `db:"username"`
	ChatID               int64     `db:"chat_id"`
	LanguageCode         *string   `db:"language"`
	IsPremium            bool      `db:"is_premium"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
	LastSeenAt           time.Time `db:"last_seen_at"`
	NotificationsEnabled bool      `db:"notifications_enabled"`
	AvatarURL            *string   `db:"avatar_url"`
	Title                *string   `db:"title"`
}

func (s *storage) getUserBy(query string, args ...interface{}) (*User, error) {
	var user User
	row := s.db.QueryRowContext(context.Background(), query, args...)

	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.LanguageCode,
		&user.ChatID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastSeenAt,
		&user.NotificationsEnabled,
		&user.AvatarURL,
		&user.Title,
	); err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *storage) GetUserByID(id int64) (*User, error) {
	return s.getUserBy("SELECT id, first_name, last_name, username, language, chat_id, created_at, updated_at, last_seen_at, notifications_enabled, avatar_url, title FROM users WHERE id = ?", id)
}

func (s *storage) GetUserByChatID(chatID int64) (*User, error) {
	return s.getUserBy("SELECT id, first_name, last_name, username, language, chat_id, created_at, updated_at, last_seen_at, notifications_enabled, avatar_url, title FROM users WHERE chat_id = ?", chatID)
}

func (s *storage) UpdateUserAvatarURL(uid int64, url string) error {
	q := `
		UPDATE users
		SET avatar_url = ?
		WHERE id = ?
	`

	res, err := s.db.Exec(q, url, uid)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return err
}

func (s *storage) DeleteUserByID(uid int64) error {
	q := `
		DELETE FROM users
		WHERE id = ?
	`

	res, err := s.db.Exec(q, uid)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return err
}

func (s *storage) UpdateUser(uid int64, user User) (*User, error) {
	q := `
		UPDATE users
		SET first_name = ?, last_name = ?, username = ?, language = ?, is_premium = ?, notifications_enabled = ?
		WHERE id = ?
	`

	res, err := s.db.Exec(q,
		user.FirstName,
		user.LastName,
		user.Username,
		user.LanguageCode,
		user.IsPremium,
		user.NotificationsEnabled,
		uid,
	)

	if err != nil {
		return nil, err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return s.GetUserByID(uid)
}

func (s *storage) CreateUser(user User) error {
	q := `
		INSERT INTO users (first_name, last_name, username, chat_id, language, is_premium, notifications_enabled, avatar_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	if _, err := s.db.Exec(q,
		user.FirstName,
		user.LastName,
		user.Username,
		user.ChatID,
		user.LanguageCode,
		user.IsPremium,
		user.NotificationsEnabled,
		user.AvatarURL,
	); err != nil {
		return err
	}

	return nil
}
