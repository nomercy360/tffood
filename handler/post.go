package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"rednit/db"
	"rednit/terrors"
	"strconv"
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
	uid := getUserID(c)

	posts, err := h.st.ListPosts(uid)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, posts)
}

type CreatePostRequest struct {
	Photo string `json:"photo" validate:"required"`
}

type UpdatePostRequest struct {
	Text     *string `json:"text"`
	Tags     []int   `json:"tags"`
	Photo    string  `json:"photo" validate:"required"`
	Location struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Address   *string  `json:"address"`
	} `json:"location"`
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
	}

	res, err := h.st.CreatePost(uid, post)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

func GetAIUpdatedPost(post *db.Post, openAIKey string) (*db.Post, error) {
	info, err := GetFoodPictureInfo(post.PhotoURL, openAIKey)
	if err != nil {
		return nil, err
	}

	post.SuggestedIngredients = info.Ingredients
	post.SuggestedDishName = &info.DishName
	post.SuggestedTags = info.Tags
	post.IsSpam = info.IsSpam

	return post, nil
}

func (h Handler) CreatePostAISuggestions(c echo.Context) error {
	uid := getUserID(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	post, err := h.runAISuggestions(uid, id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, post)
}

func (h Handler) runAISuggestions(uid, postID int64) (*db.Post, error) {
	post, err := h.st.GetPostByID(uid, postID)

	if err != nil {
		return nil, err
	}

	post, err = GetAIUpdatedPost(post, h.config.OpenAIKey)

	if err != nil {
		return nil, err
	}

	res, err := h.st.UpdatePostSuggestions(uid, postID, *post)

	if err != nil {
		return nil, err
	}

	insights, err := GetNutritionInfo(formatIngredients(post.SuggestedIngredients), h.config.OpenAIKey)

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

	if req.Location.Latitude != nil && req.Location.Longitude != nil {
		post.Location = &db.Location{
			Latitude:  *req.Location.Latitude,
			Longitude: *req.Location.Longitude,
		}

		if req.Location.Address != nil {
			post.Location.Address = *req.Location.Address
		}
	}

	res, err := h.st.UpdatePost(uid, id, post, req.Tags)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) ReactToPost(c echo.Context) error {
	uid := getUserID(c)
	postID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	reaction := db.ReactionType(c.Param("reaction"))

	if reaction != db.ReactionSmile && reaction != db.ReactionFrown && reaction != db.ReactionMeh {
		return terrors.BadRequest(fmt.Errorf("invalid reaction type: %s", reaction), "invalid data")
	}

	if err := h.st.UpdatePostReaction(uid, postID, reaction); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h Handler) DropPostReaction(c echo.Context) error {
	uid := getUserID(c)
	postID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.st.DeletePostReaction(uid, postID); err != nil && errors.Is(db.ErrNotFound, err) {
		return terrors.NotFound(err, "not found")
	} else if err != nil {
		return terrors.InternalServerError(err, "unable delete reaction")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h Handler) GetTags(c echo.Context) error {
	tags, err := h.st.ListTags()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tags)
}
