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
		    title TEXT
		);

		CREATE TABLE IF NOT EXISTS locations (
		    id INTEGER PRIMARY KEY,
		    latitude REAL NOT NULL,
		    longitude REAL NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    address TEXT,
		    user_id INTEGER NOT NULL,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS posts (
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    text TEXT,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    hidden_at TIMESTAMP,
		    photo_url TEXT,
		    ingredients TEXT,
		    dish_name TEXT,
		    is_spam BOOLEAN NOT NULL DEFAULT FALSE,
		    suggested_dish_name TEXT,
		    suggested_ingredients TEXT,
		    suggested_tags TEXT,
		    location_id INTEGER,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		    FOREIGN KEY (location_id) REFERENCES locations (id) ON DELETE SET NULL
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
		
		CREATE TABLE IF NOT EXISTS reactions (
		    user_id INTEGER NOT NULL,
		    post_id INTEGER NOT NULL,
		    type TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
		    PRIMARY KEY (user_id, post_id)
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
 	`

	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
