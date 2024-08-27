package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	telegram "github.com/go-telegram/bot"
	tgModels "github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	nanoid "github.com/matoous/go-nanoid"
	"io"
	"log"
	"net/http"
	"rednit/db"
	"regexp"
	"strconv"
	"strings"
)

var messages = map[string]map[string]string{
	"en": {
		"welcome":          "This bot will help you track your meals and get insights about your nutrition.\nTry sending a photo",
		"openWebApp":       "You can open the web app by tapping the button below.",
		"gettingInsights":  "Getting insights from the image...",
		"photoAddError":    "Please send the picture as a 'Photo', not as a 'File'.",
		"uploadError":      "Failed to upload the image. Please try again.",
		"insightsNotFound": "No insights found for this image.",
	},
	"ru": {
		"welcome":          "Этот бот поможет вам отслеживать приемы пищи и получать информацию о вашем питании.\nПопробуй отправить фото",
		"openWebApp":       "Вы можете открыть веб-приложение, нажав на кнопку ниже.",
		"gettingInsights":  "Обрабатка в процессе...",
		"photoAddError":    "Пожалуйста, отправьте изображение как 'Фото', а не как 'Файл'.",
		"uploadError":      "Не удалось загрузить изображение. Пожалуйста, попробуйте еще раз.",
		"insightsNotFound": "Для этого изображения не найдено данных.",
	},
}

func getInsightsText(lang string, insights *db.FoodInsights, dishName *string) string {
	if lang == "ru" {
		return fmt.Sprintf("*%s*\n\nКалории: %d ккал\n\nБелки: %d г\nУглеводы: %d г\nЖиры: %d г",
			*dishName, insights.Calories, insights.Proteins, insights.Carbohydrates, insights.Fats)
	} else {
		return fmt.Sprintf("*%s*\nCalories: %d kcal\n\nProteins: %d g\nСarbohydrates: %d g\nАats: %d g",
			*dishName, insights.Calories, insights.Proteins, insights.Carbohydrates, insights.Fats)
	}
}

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

		msg := messages[*user.LanguageCode]["welcome"]

		params := &telegram.SendMessageParams{ChatID: update.Message.Chat.ID, Text: msg, ReplyMarkup: &webApp, ParseMode: "Markdown"}

		if _, err := h.tg.SendMessage(context.Background(), params); err != nil {
			log.Printf("Failed to send message: %v", err)
			return c.NoContent(http.StatusOK)
		}

		go h.setMenuButton(update.Message.Chat.ID)

	} else if err != nil {
		log.Printf("Failed to get user: %v", err)
		return c.NoContent(http.StatusOK)
	} else if user != nil {
		lang := "en"
		if user.LanguageCode != nil && *user.LanguageCode == "ru" {
			lang = "ru"
		}

		log.Printf("User %d already exists, sending message", user.ChatID)

		msg := &telegram.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			ParseMode: "Markdown",
		}

		if update.Message.Photo != nil && len(update.Message.Photo) > 0 {
			if err := h.onImageMessage(lang, user.ID, update); err != nil {
				log.Printf("Failed to process image from telegram. User: %d, Error: %v", user.ID, err)
				msg.Text = messages[lang]["uploadError"]
			} else {
				msg.Text = messages[lang]["gettingInsights"]
			}

		} else if update.Message.Document != nil {
			msg.Text = messages[lang]["photoAddError"]
		} else {
			msg.Text = messages[lang]["openWebApp"]
			msg.ReplyMarkup = &webApp
		}

		if _, err := h.tg.SendMessage(context.Background(), msg); err != nil {
			log.Printf("Failed to send message: %v", err)
			return c.NoContent(http.StatusOK)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) onImageMessage(lang string, uid int64, update tgModels.Update) error {
	// find the most quality photo
	photo := update.Message.Photo[len(update.Message.Photo)-1]

	var caption *string
	if update.Message.Caption != "" {
		caption = &update.Message.Caption
	}

	key, err := h.handlePhotoUpload(uid, photo)
	if err != nil {
		log.Printf("Failed to handle photo upload: %v", err)
	}

	log.Printf("Photo uploaded to bucket: %s", *key)

	post := db.Post{
		PhotoURL: fmt.Sprintf("%s/%s", h.config.CdnURL, *key),
		Text:     caption,
	}

	res, err := h.st.CreatePost(uid, post)

	if err != nil {
		return err
	}

	// run AI model in the background
	go func() {
		postWithSuggestions, err := h.runAISuggestions(lang, uid, res.ID)
		if err != nil {
			log.Printf("Failed to run AI suggestions: %v", err)
			return
		}

		insights := postWithSuggestions.FoodInsights
		var msgText string

		if insights != nil && postWithSuggestions.DishName != nil {
			msgText = getInsightsText(lang, insights, postWithSuggestions.DishName)
		} else {
			msgText = messages[lang]["insightsNotFound"]
		}

		msg := telegram.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      msgText,
			ParseMode: "Markdown",
		}

		if _, err := h.tg.SendMessage(context.Background(), &msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}()

	return nil
}

func (h Handler) handlePhotoUpload(uid int64, photo tgModels.PhotoSize) (*string, error) {
	file, err := h.tg.GetFile(context.Background(), &telegram.GetFileParams{FileID: photo.FileID})
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		return nil, err
	}

	fileURL := h.tg.FileDownloadLink(file)

	key, err := h.handleUploadToBucket(strconv.FormatInt(uid, 10), fileURL)

	if err != nil {
		log.Printf("Download/Upload failed: %v", err)
		return nil, err
	}

	return key, nil
}

func (h Handler) handleUploadToBucket(idKey, fileURL string) (*string, error) {
	resp, err := http.Get(fileURL)

	if err != nil {
		log.Printf("Failed to download file: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download file: received non-200 response code %d", resp.StatusCode)
		return nil, err
	}

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read file content: %v", err)
		return nil, err
	}

	fileReader := bytes.NewReader(fileBytes)

	key := fmt.Sprintf("media/%s/%s.jpg", idKey, nanoid.MustID(8))

	if err := h.s3Client.UploadFile(key, fileReader); err != nil {
		log.Printf("Failed to upload file: %v", err)
		return nil, err
	}

	log.Printf("File uploaded to bucket: %s", key)

	return &key, nil
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

		fileBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read file content: %v", err)
		}

		fileReader := bytes.NewReader(fileBytes)

		fileName := fmt.Sprintf("%d/%d.jpg", userID, chatID)

		if err := h.s3Client.UploadFile(fileName, fileReader); err != nil {
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
