package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"rednit/db"
	"strconv"
	"time"
)

func (h Handler) GetPost(c echo.Context) error {
	uid := getUserID(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	user, err := h.st.GetPostByID(uid, id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h Handler) GetPosts(c echo.Context) error {
	end := time.Now()
	start := end.AddDate(0, -1, 0)

	posts, err := h.st.ListPosts(db.ListPostsParams{StartDate: start, EndDate: end})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, posts)
}

type Goal string

const (
	GainMuscles Goal = "gain_muscles"
	LoseWeight  Goal = "lose_weight"
)

func (h Handler) GetFoodInsightsHandler(c echo.Context) error {
	uid := getUserID(c)
	endString := c.QueryParam("date")

	if endString == "" {
		endString = time.Now().Format("2006-01-02T15:04:05Z")
	}

	end, err := time.Parse("2006-01-02T15:04:05Z", endString)
	if err != nil {
		return err
	}

	// 00:00:00 of the end date
	start := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	user, err := h.st.GetUserByID(db.UserQuery{ID: uid})
	if err != nil {
		return err
	}

	posts, err := h.st.ListPosts(db.ListPostsParams{UserID: &uid, StartDate: start, EndDate: end})
	if err != nil {
		return err
	}

	macros := map[string]int{
		"proteins":      0,
		"fats":          0,
		"carbohydrates": 0,
	}

	for _, post := range posts {
		if post.FoodInsights != nil {
			macros["proteins"] += post.FoodInsights.Proteins
			macros["fats"] += post.FoodInsights.Fats
			macros["carbohydrates"] += post.FoodInsights.Carbohydrates
		}
	}

	caloriesConsumed := macros["proteins"]*4 + macros["fats"]*9 + macros["carbohydrates"]*4

	genderFemale := "Female"

	// Calculate BMR based on user details
	bmr := 88.362 + (13.397 * *user.Weight) + (4.799 * float64(*user.Height)) - (5.677 * float64(*user.Age))
	if user.Gender == &genderFemale {
		bmr = 447.593 + (9.247 * *user.Weight) + (3.098 * float64(*user.Height)) - (4.330 * float64(*user.Age))
	}

	// Adjust BMR based on activity level and goals
	// For example, if activityLevel is 1.2 and user's goal is to lose weight, decrease by 15%
	activityLevel := 1.2 // This should be fetched or calculated based on user input
	calorieNeeds := bmr * activityLevel
	switch *user.Goal {
	case "gain_muscles":
		calorieNeeds *= 1.2
	case "lose_weight":
		calorieNeeds *= 0.85
	}

	caloriesLeft := int(calorieNeeds) - caloriesConsumed

	response := map[string]interface{}{
		"macros":            macros,
		"calories_left":     caloriesLeft,
		"calorie_needs":     int(calorieNeeds),
		"calories_consumed": caloriesConsumed,
	}

	return c.JSON(http.StatusOK, response)
}

type CreatePostRequest struct {
	Photo   string  `json:"photo" validate:"required"`
	Text    *string `json:"text"`
	Publish bool    `json:"publish"`
}

type UpdatePostRequest struct {
	Text  *string `json:"text"`
	Tags  []int   `json:"tags"`
	Photo string  `json:"photo" validate:"required"`
}

func (h Handler) RerunAllPostsImageRecognition(c echo.Context) error {
	now := time.Now()
	// 1 year ago
	start := now.AddDate(-1, 0, 0)
	showHidden := true

	params := db.ListPostsParams{
		StartDate:  start,
		EndDate:    now,
		ShowHidden: &showHidden,
	}

	posts, err := h.st.ListPosts(params)
	if err != nil {
		return err
	}

	for _, post := range posts {
		go func(post db.Post) {
			if _, err := h.runAISuggestions("en", post.UserID, post.ID); err != nil {
				log.Printf("Failed to run AI suggestions: %v", err)
			}
		}(post)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h Handler) CreatePost(c echo.Context) error {
	uid := getUserID(c)

	var req CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	var hiddenAt *time.Time
	if !req.Publish {
		now := time.Now()
		hiddenAt = &now
	}

	post := db.Post{
		PhotoURL: fmt.Sprintf("%s/%s", h.config.CdnURL, req.Photo),
		Text:     req.Text,
		HiddenAt: hiddenAt,
	}

	res, err := h.st.CreatePost(uid, post)

	if err != nil {
		return err
	}

	go func() {
		user, err := h.st.GetUserByID(db.UserQuery{ID: uid})
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			return
		}

		lang := "en"

		if user.LanguageCode != nil && *user.LanguageCode == "ru" {
			lang = "ru"
		}

		if _, err := h.runAISuggestions(lang, uid, res.ID); err != nil {
			log.Printf("Failed to run AI suggestions: %v", err)
		}
	}()

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) runAISuggestions(lang string, uid, postID int64) (*db.Post, error) {
	post, err := h.st.GetPostByID(uid, postID)

	if err != nil {
		return nil, err
	}

	info, err := getFoodPictureInfo(lang, post.PhotoURL, post.Text, h.config.OpenAIKey)
	if err != nil {
		return nil, err
	}

	if info.IsSpam {
		if err := h.st.MarkPostAsSpam(uid, postID); err != nil {
			return nil, err
		}

		post.IsSpam = true
		return post, nil
	}

	insights, err := getNutritionInfo(lang, formatIngredients(lang, info.Ingredients), h.config.OpenAIKey)

	if err != nil {
		return nil, err
	}

	post.DishName = &info.DishName
	post.Ingredients = insights.Ingredients

	var protein, fats, carbohydrates, calories int
	for _, ingredient := range insights.Ingredients {
		protein += int(ingredient.Macros.Proteins)
		fats += int(ingredient.Macros.Fats)
		carbohydrates += int(ingredient.Macros.Carbohydrates)
		calories += int(ingredient.Calories)
	}

	post.FoodInsights = &db.FoodInsights{
		Calories:      calories,
		Proteins:      protein,
		Fats:          fats,
		Carbohydrates: carbohydrates,
	}

	res, err := h.st.UpdatePost(uid, postID, *post, nil)

	if err != nil {
		return nil, err
	}

	return res, err
}

func (h Handler) UpdatePost(c echo.Context) error {
	uid := getUserID(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	post := db.Post{
		Text:     req.Text,
		PhotoURL: fmt.Sprintf("%s/%s", h.config.CdnURL, req.Photo),
	}

	res, err := h.st.UpdatePost(uid, id, post, req.Tags)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) GetTags(c echo.Context) error {
	tags, err := h.st.ListTags()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tags)
}
