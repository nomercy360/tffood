package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	telegram "github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"rednit/db"
	"regexp"
	"strconv"
	"strings"
)

func (h Handler) HandleWebhook(c echo.Context) error {
	var update tgModels.Update
	if err := json.NewDecoder(c.Request().Body).Decode(&update); err != nil {
		return c.String(http.StatusBadRequest, "Invalid update")
	}

	if update.Message == nil {
		return c.String(http.StatusBadRequest, "No message")
	}

	if update.Message.Chat.Type != "private" {
		return c.NoContent(http.StatusOK)
	} else if update.Message.From.IsBot {
		return c.NoContent(http.StatusOK)
	}

	user, err := h.st.GetUserByID(db.UserQuery{ChatID: update.Message.Chat.ID})

	if update.Message.Text == "/reset" && user != nil {
		msg := telegram.SendMessageParams{
			ChatID: update.Message.Chat.ID,
		}

		if err := h.st.DeleteUserByID(user.ID); err != nil {
			log.Printf("Failed to delete user: %v", err)
			msg.Text = "Failed to delete user"
		} else {
			msg.Text = "User deleted"
		}

		if _, err := h.tg.SendMessage(context.Background(), &msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}

		return c.NoContent(http.StatusOK)
	}

	webApp := tgModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgModels.InlineKeyboardButton{
			{
				{Text: "start", WebApp: &tgModels.WebAppInfo{URL: h.config.WebAppURL}},
			},
		},
	}

	if err != nil && errors.Is(err, db.ErrNotFound) {
		log.Printf("User %d not found, creating new user", update.Message.Chat.ID)

		user = h.createUser(update)
		if user == nil {
			return c.NoContent(http.StatusOK)
		}

		photo := &tgModels.InputFileString{Data: "https://assets.peatch.io/peatch-preview.png"}

		params := &telegram.SendPhotoParams{ChatID: update.Message.Chat.ID, Caption: "Hello", ReplyMarkup: &webApp, Photo: photo, ParseMode: "Markdown"}

		if _, err := h.tg.SendPhoto(context.Background(), params); err != nil {
			log.Printf("Failed to send message: %v", err)
			return c.NoContent(http.StatusOK)
		}

		go h.setMenuButton(update.Message.Chat.ID)

	} else if err != nil {
		log.Printf("Failed to get user: %v", err)
		return c.NoContent(http.StatusOK)
	} else {
		log.Printf("User %d already exists, sending message", user.ChatID)

		params := &telegram.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Open App", ReplyMarkup: &webApp, ParseMode: "Markdown"}

		if _, err := h.tg.SendMessage(context.Background(), params); err != nil {
			log.Printf("Failed to send message: %v", err)
			return c.NoContent(http.StatusOK)
		}
	}

	return c.NoContent(http.StatusOK)
}

func extractReferrerID(arg string) int64 {
	idStr := strings.TrimPrefix(arg, "friend")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Failed to parse referrer ID: %v", err)
		return 0
	}
	return id
}

func (h Handler) createUser(update tgModels.Update) *db.User {
	// Extract user details from update
	var firstName, lastName *string
	if update.Message.Chat.FirstName != "" {
		firstName = &update.Message.Chat.FirstName
	}

	if update.Message.Chat.LastName != "" {
		lastName = &update.Message.Chat.LastName
	}

	// if username is empty, use first name
	username := update.Message.Chat.Username

	if username == "" {
		username = "user_" + fmt.Sprintf("%d", update.Message.Chat.ID)
	}

	user := db.User{
		ChatID:    update.Message.Chat.ID,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	}

	lang := "ru"

	if update.Message.From.LanguageCode != "ru" {
		lang = "en"
	}

	user.LanguageCode = &lang

	newUser, err := h.st.CreateUser(user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil
	}

	go h.handleUserAvatar(newUser.ID, update.Message.From.ID, newUser.ChatID)

	return newUser
}

func (h Handler) handleUserAvatar(userID, tgUserID, chatID int64) {
	photos, err := h.tg.GetUserProfilePhotos(context.Background(), &telegram.GetUserProfilePhotosParams{UserID: tgUserID, Offset: 0, Limit: 1})
	if err != nil {
		log.Printf("Failed to get user profile photos: %v", err)
		return
	}

	if photos.TotalCount > 0 {
		bestPhoto := new(tgModels.PhotoSize)

		for _, album := range photos.Photos {
			for _, pic := range album {
				if pic.FileSize > bestPhoto.FileSize || (pic.FileSize == bestPhoto.FileSize && pic.Width > bestPhoto.Width) {
					bestPhoto = &pic
				}
			}
		}

		file, err := h.tg.GetFile(context.Background(), &telegram.GetFileParams{FileID: bestPhoto.FileID})
		if err != nil {
			log.Printf("Failed to get file: %v", err)
			return
		}

		fileURL := h.tg.FileDownloadLink(file)

		resp, err := http.Get(fileURL)

		if err != nil {
			log.Printf("Failed to download file: %v", err)
			return
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Failed to read file: %v", err)
			return
		}

		fileName := fmt.Sprintf("%d/%d.jpg", userID, chatID)

		if err := h.s3Client.UploadFile(data, fileName); err != nil {
			log.Printf("Failed to upload user avatar to S3: %v", err)
			return
		}

		log.Printf("Avatar uploaded successfully: %s", fileName)

		url := fmt.Sprintf("%s/%s", h.config.CdnURL, fileName)

		if err := h.st.UpdateUserAvatarURL(userID, url); err != nil {
			log.Printf("Failed to update user avatar URL: %v", err)
		}

		log.Printf("Profile photo updated for user %d", chatID)
	}
}

func (h Handler) setMenuButton(chatID int64) {
	menu := telegram.SetChatMenuButtonParams{
		ChatID: chatID,
		MenuButton: tgModels.MenuButtonWebApp{
			Type:   "web_app",
			Text:   "Open App",
			WebApp: tgModels.WebAppInfo{URL: h.config.WebAppURL},
		},
	}

	if _, err := h.tg.SetChatMenuButton(context.Background(), &menu); err != nil {
		log.Printf("Failed to set chat menu button: %v", err)
		return
	}

	log.Printf("User %d menu button set", chatID)
}

func urlify(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, " ", "_")

	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	s = reg.ReplaceAllString(s, "_")

	s = strings.Trim(s, "_")

	return s
}
