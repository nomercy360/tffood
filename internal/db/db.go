package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"time"
)

type storage struct {
	db *sql.DB
}

func init() {
	// Registers the sqlite3 driver with a ConnectHook so that we can
	// initialize the default PRAGMAs.
	//
	// Note 1: we don't define the PRAGMA as part of the dsn string
	// because not all pragmas are available.
	//
	// Note 2: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	sql.Register("sql",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				_, err := conn.Exec(`
					PRAGMA busy_timeout       = 10000;
					PRAGMA journal_mode       = WAL;
					PRAGMA journal_size_limit = 200000000;
					PRAGMA synchronous        = NORMAL;
					PRAGMA foreign_keys       = ON;
					PRAGMA temp_store         = MEMORY;
					PRAGMA cache_size         = -16000;
				`, nil)

				return err
			},
		},
	)
}

func NewStorage(dbFile string) (*storage, error) {
	db, err := sql.Open("sql", dbFile)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	createTables := `
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
		    UNIQUE (chat_id)
		);

		CREATE TABLE IF NOT EXISTS meals (
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
		    aesthetic_rating INTEGER,
		    health_rating INTEGER,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS comments (
		    id INTEGER PRIMARY KEY,
		    user_id INTEGER NOT NULL,
		    meal_id INTEGER NOT NULL,
		    text TEXT NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		    FOREIGN KEY (meal_id) REFERENCES meals (id) ON DELETE CASCADE
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

		CREATE TABLE IF NOT EXISTS meal_tags (
		    meal_id INTEGER NOT NULL,
		    tag_id INTEGER NOT NULL,
		    PRIMARY KEY (meal_id, tag_id),
		    FOREIGN KEY (meal_id) REFERENCES meals (id) ON DELETE CASCADE,
		    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
		);

		INSERT INTO tags (name) VALUES 
		('Keto'), ('Breakfast'), ('Lunch'), ('Dinner'), ('Snack'), ('Vegetarian'), ('Vegan')
		ON CONFLICT DO NOTHING;
	`
	_, err = db.Exec(createTables)
	if err != nil {
		return nil, err
	}

	return &storage{db: db}, nil
}

type HealthStats struct {
	Status            string `json:"status"`
	Error             string `json:"error,omitempty"`
	Message           string `json:"message"`
	OpenConnections   int    `json:"open_connections"`
	InUse             int    `json:"in_use"`
	Idle              int    `json:"idle"`
	WaitCount         int64  `json:"wait_count"`
	WaitDuration      string `json:"wait_duration"`
	MaxIdleClosed     int64  `json:"max_idle_closed"`
	MaxLifetimeClosed int64  `json:"max_lifetime_closed"`
}

func (s *storage) Health() (HealthStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := HealthStats{}

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats.Status = "down"
		stats.Error = fmt.Sprintf("db down: %v", err)
		return stats, fmt.Errorf("db down: %w", err)
	}

	// Database is up, add more statistics
	stats.Status = "up"
	stats.Message = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats.OpenConnections = dbStats.OpenConnections
	stats.InUse = dbStats.InUse
	stats.Idle = dbStats.Idle
	stats.WaitCount = dbStats.WaitCount
	stats.WaitDuration = dbStats.WaitDuration.String()
	stats.MaxIdleClosed = dbStats.MaxIdleClosed
	stats.MaxLifetimeClosed = dbStats.MaxLifetimeClosed

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats.Message = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats.Message = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats.Message = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats, nil
}
