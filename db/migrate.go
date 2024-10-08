package db

func (s Storage) Migrate() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
		    id INTEGER PRIMARY KEY,
		    username TEXT NOT NULL,
		    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    language TEXT NOT NULL DEFAULT 'en',
		    first_name TEXT,
		    last_name TEXT,
		    chat_id INTEGER NOT NULL,
		    last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    notifications_enabled BOOLEAN NOT NULL DEFAULT TRUE,
		    avatar_url TEXT,
		    title TEXT,
		    age INTEGER,
		    weight DECIMAL(5, 2),
		    height INTEGER,
		    fat_percentage DECIMAL(5, 2),
		    activity_level TEXT,
		    goal TEXT,
		    gender TEXT,
		    UNIQUE (chat_id)
		);

		CREATE TABLE IF NOT EXISTS posts (
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    text TEXT,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    hidden_at TIMESTAMP,
		    photo_url TEXT,
		    is_spam BOOLEAN NOT NULL DEFAULT FALSE,
		    dish_name TEXT,
		    ingredients TEXT,
		    tags TEXT,
		    food_insights TEXT,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS comments (
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    post_id INTEGER NOT NULL,
		    text TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS followers (
		    id INTEGER PRIMARY KEY,
		    follower_id INTEGER NOT NULL,
		    followee_id INTEGER NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
		    FOREIGN KEY (followee_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS tags (
		    id INTEGER PRIMARY KEY,
		    name TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    language TEXT NOT NULL DEFAULT 'en',
		    UNIQUE (name)       
		);

		CREATE TABLE IF NOT EXISTS post_tags (
		    post_id INTEGER NOT NULL,
		    tag_id INTEGER NOT NULL,
		    PRIMARY KEY (post_id, tag_id),
		    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
		    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
		);

		INSERT INTO tags (name) VALUES 
		('Keto'), ('Breakfast'), ('Lunch'), ('Dinner'), ('Snack'), ('Vegetarian'), ('Vegan')
		ON CONFLICT DO NOTHING;

		CREATE TABLE IF NOT EXISTS bot_messages (
			id INTEGER PRIMARY KEY,
			chat_id INTEGER NOT NULL,
			message_id INTEGER NOT NULL,
			entity_id INTEGER NOT NULL,
			sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
 	`

	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
