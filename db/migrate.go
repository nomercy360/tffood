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

		CREATE TABLE IF NOT EXISTS posts (
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    text TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    hidden_at TIMESTAMP,
		    photo_url TEXT,
		    FOREIGN KEY (user_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS comments (
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    post_id INTEGER NOT NULL,
		    text TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users (id),
		    FOREIGN KEY (post_id) REFERENCES posts (id)
		);
		
		CREATE TABLE IF NOT EXISTS reactions (
		    user_id INTEGER NOT NULL,
		    post_id INTEGER NOT NULL,
		    type TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users (id),
		    FOREIGN KEY (post_id) REFERENCES posts (id),
		    PRIMARY KEY (user_id, post_id)
		);

		CREATE TABLE IF NOT EXISTS followers (
		    id INTEGER PRIMARY KEY,
		    follower_id INTEGER NOT NULL,
		    followee_id INTEGER NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (follower_id) REFERENCES users (id),
		    FOREIGN KEY (followee_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS locations (
		    id INTEGER PRIMARY KEY,
		    latitude REAL NOT NULL,
		    longitude REAL NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    name TEXT
		);

		CREATE TABLE IF NOT EXISTS posts_locations (
		    post_id INTEGER NOT NULL,
		    location_id INTEGER NOT NULL,
		    PRIMARY KEY (post_id, location_id),
		    FOREIGN KEY (post_id) REFERENCES posts (id),
		    FOREIGN KEY (location_id) REFERENCES locations (id)
		);
		
		CREATE TABLE IF NOT EXISTS tags (
		    id INTEGER PRIMARY KEY,
		    name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS post_tags (
		    post_id INTEGER NOT NULL,
		    tag_id INTEGER NOT NULL,
		    PRIMARY KEY (post_id, tag_id),
		    FOREIGN KEY (post_id) REFERENCES posts (id),
		    FOREIGN KEY (tag_id) REFERENCES tags (id)
		);
 	`

	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	return nil
}
