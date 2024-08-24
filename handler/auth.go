package handler

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/telegram-mini-apps/init-data-golang"
	"math/rand"
	"net/http"
	"rednit/db"
	"rednit/terrors"
	"time"
)

type UserWithToken struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

func getUserID(c echo.Context) int64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return claims.UID
}

func (h Handler) TelegramAuth(c echo.Context) error {
	expIn := 24 * time.Hour
	botToken := h.config.BotToken

	query := c.QueryString()

	if err := initdata.Validate(query, botToken, expIn); err != nil {
		return terrors.Unauthorized(err, "invalid data")
	}

	data, err := initdata.Parse(query)

	if err != nil {
		return terrors.Unauthorized(err, "invalid data")
	}

	user, err := h.st.GetUserByID(db.UserQuery{ChatID: data.User.ID})

	if err != nil && errors.Is(err, db.ErrNotFound) {
		var firstName, lastName *string

		if data.User.FirstName != "" {
			firstName = &data.User.FirstName
		}

		if data.User.LastName != "" {
			lastName = &data.User.LastName
		}

		username := data.User.Username
		if username == "" {
			username = "u_" + fmt.Sprintf("%d", data.User.ID)
		}

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		randNum := r1.Intn(39) + 1

		profilePic := fmt.Sprintf("https://fm-assets.mxksim.dev/avatars/%d.svg", randNum)

		create := db.User{
			FirstName:            firstName,
			LastName:             lastName,
			Username:             username,
			ChatID:               data.User.ID,
			IsPremium:            data.User.IsPremium,
			NotificationsEnabled: data.User.AllowsWriteToPm,
			AvatarURL:            &profilePic,
		}

		lang := "ru"

		if data.User.LanguageCode != "ru" {
			lang = "en"
		}

		create.LanguageCode = &lang

		user, err = h.st.CreateUser(create)
		if err != nil {
			return terrors.InternalServerError(err, "failed to create user")
		}
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get user")
	}

	token, err := generateJWT(h.config.JWTSecret, user.ID, user.ChatID)

	if err != nil {
		return terrors.InternalServerError(err, "failed to generate token")
	}

	return c.JSON(http.StatusOK, UserWithToken{
		User:  *user,
		Token: token,
	})
}

type UpdateUserRequest struct {
	Language             string `json:"language"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
}

func (h Handler) UpdateUserSettings(c echo.Context) error {
	uid := getUserID(c)

	var req UpdateUserRequest

	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to bind request")
	}

	user, err := h.st.GetUserByID(db.UserQuery{ID: uid})

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "user not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to get user")
	}

	user.LanguageCode = &req.Language
	user.NotificationsEnabled = req.NotificationsEnabled

	updated, err := h.st.UpdateUser(uid, *user)

	if err != nil {
		return terrors.InternalServerError(err, "failed to update user")
	}

	return c.JSON(http.StatusOK, updated)
}

func (h Handler) SubmitJoinCommunityRequest(c echo.Context) error {
	uid := getUserID(c)

	err := h.st.UpdateUserRequestToJoin(uid)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "failed to update user")
	}

	return c.NoContent(http.StatusNoContent)
}
