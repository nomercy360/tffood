package api

import (
	"eatsome/internal/db"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getUserID(c echo.Context) int64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return claims.UID
}

type UserResponse struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type MealResponse struct {
	ID              int64            `json:"id"`
	UserID          int64            `json:"user_id"`
	PhotoURL        string           `json:"photo_url"`
	Text            *string          `json:"text"`
	DishName        *string          `json:"dish_name"`
	AestheticRating *int             `json:"aesthetic_rating"`
	HealthRating    *int             `json:"health_rating"`
	IsSpam          bool             `json:"is_spam"`
	FoodInsights    *db.FoodInsights `json:"food_insights"`
	User            UserResponse     `json:"user"`
	Ingredients     db.Ingredients   `json:"ingredients"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func (a *API) GetMeals(c echo.Context) error {
	end := time.Now()
	start := end.AddDate(0, -1, 0)

	meals, err := a.storage.ListMeals(start, end)

	if err != nil {
		return err
	}

	var resp []MealResponse

	// fetch user data
	for _, meal := range meals {
		user, err := a.storage.GetUserByID(meal.UserID)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			continue
		}

		resp = append(resp, MealResponse{
			ID:              meal.ID,
			UserID:          meal.UserID,
			PhotoURL:        meal.PhotoURL,
			Text:            meal.Text,
			DishName:        meal.DishName,
			AestheticRating: meal.AestheticRating,
			HealthRating:    meal.HealthRating,
			IsSpam:          meal.IsSpam,
			FoodInsights:    meal.FoodInsights,
			Ingredients:     meal.Ingredients,
			CreatedAt:       meal.CreatedAt,
			UpdatedAt:       meal.UpdatedAt,
			User: UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				AvatarURL: user.AvatarURL,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
		})
	}

	return c.JSON(http.StatusOK, resp)
}

type CreateMealRequest struct {
	Photo string  `json:"photo" validate:"required"`
	Text  *string `json:"text"`
}

type UpdateMealRequest struct {
	Text  *string `json:"text"`
	Tags  []int   `json:"tags"`
	Photo string  `json:"photo" validate:"required"`
}

func (a *API) CreateMeal(c echo.Context) error {
	uid := getUserID(c)

	var req CreateMealRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	meal := db.Meal{
		PhotoURL: fmt.Sprintf("%s/%s", a.cfg.AssetsURL, req.Photo),
		Text:     req.Text,
	}

	res, err := a.storage.AddMeal(uid, meal)

	if err != nil {
		return err
	}

	go func() {
		user, err := a.storage.GetUserByID(uid)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			return
		}

		lang := "en"

		if user.LanguageCode != nil && *user.LanguageCode == "ru" {
			lang = "ru"
		}

		if _, err := a.runAISuggestions(lang, uid, res.ID); err != nil {
			log.Printf("Failed to run AI suggestions: %v", err)
		}
	}()

	return c.JSON(http.StatusCreated, res)
}

func (a *API) runAISuggestions(lang string, uid, mealID int64) (*db.Meal, error) {
	meal, err := a.storage.GetMealByID(mealID)

	if err != nil {
		return nil, err
	}

	info, err := a.recognizer.GetFoodPictureInfo(lang, meal.PhotoURL, meal.Text)
	if err != nil {
		return nil, err
	}

	meal.IsSpam = info.IsSpam
	meal.DishName = &info.DishName
	meal.AestheticRating = &info.AestheticRating
	meal.HealthRating = &info.HealthRating

	meal.FoodInsights = &db.FoodInsights{
		Calories:      info.Calories,
		Proteins:      info.Proteins,
		Fats:          info.Fats,
		Carbohydrates: info.Carbohydrates,
	}

	meal.Ingredients = info.IngredientsInfo

	res, err := a.storage.UpdateMeal(uid, mealID, *meal, nil)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (a *API) UpdateMeal(c echo.Context) error {
	uid := getUserID(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req UpdateMealRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	meal := db.Meal{
		Text:     req.Text,
		PhotoURL: fmt.Sprintf("%s/%s", a.cfg.AssetsURL, req.Photo),
	}

	res, err := a.storage.UpdateMeal(uid, id, meal, req.Tags)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}
