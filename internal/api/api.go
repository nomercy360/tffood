package api

import (
	"eatsome/internal/db"
	"eatsome/internal/recognition"
	"eatsome/internal/s3"
	"time"
)

// storager interface for database operations
type storager interface {
	Health() (db.HealthStats, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	GetUserByID(id int64) (*db.User, error)
	CreateUser(user db.User) error
	GetMealByID(id int64) (*db.Meal, error)
	ListMeals(startDate, endDate time.Time) ([]db.Meal, error)
	AddMeal(uid int64, meal db.Meal) (*db.Meal, error)
	UpdateMeal(uid, id int64, meal db.Meal, tags []int) (*db.Meal, error)
}

type API struct {
	storage  storager
	s3Client *s3.Client

	recognizer *recognition.Client

	// Config struct
	cfg Config
}

type Config struct {
	BotToken  string
	JWTSecret string
	AssetsURL string
}

func New(storage storager, cfg Config, s3Client *s3.Client, recognizer *recognition.Client) *API {
	return &API{
		storage:    storage,
		cfg:        cfg,
		s3Client:   s3Client,
		recognizer: recognizer,
	}
}
