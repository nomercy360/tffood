package handler

import (
	telegram "github.com/go-telegram/bot"
	"io"
	"rednit/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"rednit/db"
)

type s3Client interface {
	GetPresignedURL(fileName string, exp time.Duration) (string, error)
	UploadFile(key string, file io.Reader) error
}

type storage interface {
	GetUserByID(params db.UserQuery) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
	CreatePost(uid int64, post db.Post) (*db.Post, error)
	UpdatePost(uid, postID int64, post db.Post, tags []int) (*db.Post, error)
	ListPosts(uid *int64, startDate, endDate time.Time) ([]db.Post, error)
	GetPostByID(uid, id int64) (*db.Post, error)
	ListTags() ([]db.Tag, error)
	UpdateUserAvatarURL(uid int64, url string) error
	DeleteUserByID(uid int64) error
	UpdateUser(uid int64, user db.User) (*db.User, error)
	UpdateUserRequestToJoin(uid int64) error
}

type Handler struct {
	st       storage
	config   config.Default
	s3Client s3Client
	tg       *telegram.Bot
}

func New(st storage, config config.Default, s3Client s3Client, tg *telegram.Bot) Handler {
	return Handler{st: st, config: config, s3Client: s3Client, tg: tg}
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    int64 `json:"uid"`
	ChatID int64 `json:"chat_id"`
}

func generateJWT(secret string, uid, chatID int64) (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    uid,
		ChatID: chatID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
