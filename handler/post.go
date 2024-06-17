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
	Text     string `json:"text"`
	Photo    string `json:"photo" validate:"required"`
	Tags     []int  `json:"tags"`
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

	info, err := GetFoodPictureInfo(post.PhotoURL, h.config.OpenAIKey)

	if err != nil {
		return err
	}

	post.SuggestedIngredients = info.Ingredients
	post.SuggestedDishName = &info.DishName
	post.IsSpam = info.IsSpam

	res, err := h.st.CreatePost(uid, post, req.Tags)

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
