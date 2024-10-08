package db

import "time"

type User struct {
	ID                   int64     `db:"id" json:"id"`
	FirstName            *string   `db:"first_name" json:"first_name"`
	LastName             *string   `db:"last_name" json:"last_name"`
	Username             string    `db:"username" json:"username"`
	ChatID               int64     `db:"chat_id" json:"chat_id"`
	LanguageCode         *string   `db:"language" json:"language"`
	IsPremium            bool      `db:"is_premium" json:"is_premium"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
	LastSeenAt           time.Time `db:"last_seen_at" json:"last_seen_at"`
	NotificationsEnabled bool      `db:"notifications_enabled" json:"notifications_enabled"`
	AvatarURL            *string   `db:"avatar_url" json:"avatar_url"`
	Title                *string   `db:"title" json:"title"`
	Age                  *int      `db:"age" json:"age"`
	Weight               *float64  `db:"weight" json:"weight"`
	Height               *int      `db:"height" json:"height"`
	FatPercentage        *float64  `db:"fat_percentage" json:"fat_percentage"`
	Goal                 *string   `db:"goal" json:"goal"`
	Gender               *string   `db:"gender" json:"gender"`
}

func (u User) GetUserLanguage() string {
	lang := "en"
	if u.LanguageCode != nil && *u.LanguageCode == "ru" {
		lang = "ru"
	}

	return lang
}

type UserQuery struct {
	ChatID int64
	ID     int64
}

func (s Storage) GetUserByID(query UserQuery) (*User, error) {
	var user User

	var args interface{}

	q := `
		SELECT id, first_name, last_name, username, chat_id, language, is_premium, created_at, updated_at, last_seen_at, notifications_enabled, avatar_url,
			title, age, weight, height, fat_percentage, goal, gender
		FROM users
	`

	if query.ChatID != 0 {
		q += " WHERE chat_id = ?"
		args = query.ChatID
	} else {
		q += " WHERE id = ?"
		args = query.ID
	}

	err := s.db.QueryRow(q, args).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.ChatID,
		&user.LanguageCode,
		&user.IsPremium,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastSeenAt,
		&user.NotificationsEnabled,
		&user.AvatarURL,
		&user.Title,
		&user.Age,
		&user.Weight,
		&user.Height,
		&user.FatPercentage,
		&user.Goal,
		&user.Gender,
	)

	if IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s Storage) CreateUser(user User) (*User, error) {
	q := `
		INSERT INTO users (first_name, last_name, username, chat_id, language, is_premium, notifications_enabled, avatar_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := s.db.Exec(q,
		user.FirstName,
		user.LastName,
		user.Username,
		user.ChatID,
		user.LanguageCode,
		user.IsPremium,
		user.NotificationsEnabled,
		user.AvatarURL,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(UserQuery{ID: id})
}

func (s Storage) UpdateUserAvatarURL(uid int64, url string) error {
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

func (s Storage) DeleteUserByID(uid int64) error {
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

func (s Storage) UpdateUser(uid int64, user User) (*User, error) {
	q := `
		UPDATE users
		SET first_name = ?, last_name = ?, username = ?, language = ?,
		    is_premium = ?, notifications_enabled = ?,
		    age = ?, weight = ?, height = ?, fat_percentage = ?, goal = ?, gender = ?
		WHERE id = ?
	`

	res, err := s.db.Exec(q,
		user.FirstName,
		user.LastName,
		user.Username,
		user.LanguageCode,
		user.IsPremium,
		user.NotificationsEnabled,
		user.Age,
		user.Weight,
		user.Height,
		user.FatPercentage,
		user.Goal,
		user.Gender,
		uid,
	)

	if err != nil {
		return nil, err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return s.GetUserByID(UserQuery{ID: uid})
}
