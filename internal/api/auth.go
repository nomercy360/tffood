package api

import (
	"eatsome/internal/contract"
	"eatsome/internal/db"
	"eatsome/internal/terrors"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateReferralCode(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

func (a *API) AuthTelegram(c echo.Context) error {
	query := c.QueryString()

	expIn := 24 * time.Hour
	botToken := a.cfg.BotToken

	if err := initdata.Validate(query, botToken, expIn); err != nil {
		return terrors.Unauthorized(err, "invalid init data from telegram")
	}

	data, err := initdata.Parse(query)

	if err != nil {
		return terrors.Unauthorized(err, "cannot parse init data from telegram")
	}

	user, err := a.storage.GetUserByChatID(data.User.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		username := data.User.Username
		if username == "" {
			username = "user_" + fmt.Sprintf("%d", data.User.ID)
		}

		var first, last *string

		if data.User.FirstName != "" {
			first = &data.User.FirstName
		}

		if data.User.LastName != "" {
			last = &data.User.LastName
		}

		lang := "ru"

		if data.User.LanguageCode != "ru" {
			lang = "en"
		}

		imgUrl := fmt.Sprintf("%s/avatars/%d.svg", a.cfg.AssetsURL, rand.Intn(30)+1)

		if data.User.PhotoURL != "" {
			imgFile := fmt.Sprintf("fb/users/%s.jpg", gonanoid.Must(8))
			imgUrl = fmt.Sprintf("%s/%s", a.cfg.AssetsURL, imgFile)
			if err = a.uploadImageToS3(data.User.PhotoURL, imgFile); err != nil {
				return terrors.InternalServerError(err, "cannot upload user avatar to S3")
			}
		}

		create := db.User{
			Username:     username,
			ChatID:       data.User.ID,
			FirstName:    first,
			LastName:     last,
			AvatarURL:    &imgUrl,
			LanguageCode: &lang,
		}

		if err = a.storage.CreateUser(create); err != nil {
			return terrors.InternalServerError(err, "cannot create user")
		}

		user, err = a.storage.GetUserByChatID(data.User.ID)
		if err != nil {
			return terrors.InternalServerError(err, "cannot get user")
		}
	} else if err != nil {
		return terrors.InternalServerError(err, "cannot get user")
	}

	token, err := generateJWT(user.ID, user.ChatID, a.cfg.JWTSecret)

	if err != nil {
		return terrors.InternalServerError(err, "jwt library error")
	}

	userResp := contract.UserResponse{
		ID:                   user.ID,
		Username:             user.Username,
		LanguageCode:         user.LanguageCode,
		ChatID:               user.ChatID,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            user.UpdatedAt,
		LastSeenAt:           user.LastSeenAt,
		NotificationsEnabled: user.NotificationsEnabled,
		AvatarURL:            user.AvatarURL,
		Title:                user.Title,
	}

	resp := &contract.UserAuthResponse{
		Token: token,
		User:  userResp,
	}

	return c.JSON(http.StatusOK, resp)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    int64 `json:"uid"`
	ChatID int64 `json:"chat_id"`
}

func generateJWT(id int64, chatID int64, secretKey string) (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    id,
		ChatID: chatID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (a *API) uploadImageToS3(imgURL string, fileName string) error {
	resp, err := http.Get(imgURL)

	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)

	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if _, err = a.s3Client.UploadFile(data, fileName); err != nil {
		return fmt.Errorf("failed to upload user avatar to S3: %v", err)
	}

	return nil
}
