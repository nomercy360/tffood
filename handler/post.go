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

	posts, err := h.st.ListPosts(nil, start, end)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, posts)
}

func (h Handler) GetFoodInsightsHandler(c echo.Context) error {
	uid := getUserID(c)
	end := time.Now()
	// one week ago
	start := end.AddDate(0, 0, -7)

	posts, err := h.st.ListPosts(&uid, start, end)
	if err != nil {
		return err
	}

	dayOrder := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	caloricBreakdown := map[string]int{
		"Mon": 0, "Tue": 0, "Wed": 0, "Thu": 0, "Fri": 0, "Sat": 0, "Sun": 0,
	}

	macros := map[string]int{
		"proteins":      0,
		"fats":          0,
		"carbohydrates": 0,
	}

	for _, post := range posts {
		day := post.CreatedAt.Weekday().String()[:3] // Get short weekday name
		if post.FoodInsights != nil {
			caloricBreakdown[day] += post.FoodInsights.Calories
			macros["proteins"] += post.FoodInsights.Proteins
			macros["fats"] += post.FoodInsights.Fats
			macros["carbohydrates"] += post.FoodInsights.Carbohydrates
		}
	}

	orderedCaloricBreakdown := make([]int, len(dayOrder))
	for i, day := range dayOrder {
		orderedCaloricBreakdown[i] = caloricBreakdown[day]
	}

	response := map[string]interface{}{
		"caloric_breakdown": orderedCaloricBreakdown,
		"macros":            macros,
	}

	return c.JSON(http.StatusOK, response)
}

type CreatePostRequest struct {
	Photo string  `json:"photo" validate:"required"`
	Text  *string `json:"text"`
}

type UpdatePostRequest struct {
	Text  *string `json:"text"`
	Tags  []int   `json:"tags"`
	Photo string  `json:"photo" validate:"required"`
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

	post := db.Post{
		PhotoURL: fmt.Sprintf("%s/%s", h.config.CdnURL, req.Photo),
		Text:     req.Text,
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

func getAIUpdatedPost(lang string, post *db.Post, openAIKey string) (*db.Post, error) {
	info, err := getFoodPictureInfo(lang, post.PhotoURL, post.Text, openAIKey)
	if err != nil {
		return nil, err
	}

	post.SuggestedIngredients = info.Ingredients
	post.SuggestedDishName = &info.DishName
	post.SuggestedTags = info.Tags
	post.IsSpam = info.IsSpam

	return post, nil
}

func (h Handler) runAISuggestions(lang string, uid, postID int64) (*db.Post, error) {
	post, err := h.st.GetPostByID(uid, postID)

	if err != nil {
		return nil, err
	}

	post, err = getAIUpdatedPost(lang, post, h.config.OpenAIKey)

	if err != nil {
		return nil, err
	}

	res, err := h.st.UpdatePostSuggestions(uid, postID, *post)

	if err != nil {
		return nil, err
	}

	insights, err := getNutritionInfo(lang, formatIngredients(post.SuggestedIngredients), h.config.OpenAIKey)

	if err != nil {
		return nil, err
	}

	fi := db.FoodInsights{
		Calories:           int(insights.Calories),
		Carbohydrates:      int(insights.Macros.Carbs),
		Fats:               int(insights.Macros.Fats),
		Proteins:           int(insights.Macros.Proteins),
		DietaryInformation: insights.DietaryInfo,
	}

	res, err = h.st.UpdatePostFoodInsights(uid, postID, fi)

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
