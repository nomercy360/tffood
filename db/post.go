package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int64   `db:"id" json:"id"`
	Username  string  `db:"username" json:"username"`
	AvatarURL string  `db:"avatar_url" json:"avatar_url"`
	FirstName *string `db:"first_name" json:"first_name"`
	LastName  *string `db:"last_name" json:"last_name"`
	Title     *string `db:"title" json:"title"`
}

func (up *UserProfile) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}

	if len(source) == 0 {
		source = []byte("{}")
	}

	if err := json.Unmarshal(source, up); err != nil {
		return fmt.Errorf("error unmarshalling UserProfile JSON: %w", err)
	}

	return nil
}

type Post struct {
	ID                   int64         `db:"id" json:"id"`
	UserID               int64         `db:"user_id" json:"user_id"`
	Text                 *string       `db:"text" json:"text"`
	CreatedAt            time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time     `db:"updated_at" json:"updated_at"`
	HiddenAt             *time.Time    `db:"hidden_at" json:"hidden_at"`
	PhotoURL             string        `db:"photo_url" json:"photo_url"`
	User                 UserProfile   `json:"user" db:"user"`
	SuggestedDishName    *string       `json:"suggested_dish_name" db:"suggested_dish_name"`
	SuggestedIngredients Ingredients   `json:"suggested_ingredients" db:"suggested_ingredients"`
	SuggestedTags        ArrayString   `json:"suggested_tags" db:"suggested_tags"`
	IsSpam               bool          `json:"is_spam" db:"is_spam"`
	Tags                 TagSlice      `json:"tags" db:"tags"`
	FoodInsights         *FoodInsights `json:"food_insights" db:"food_insights"`
}

type FoodInsights struct {
	Calories           int      `json:"calories" db:"calories"`
	Proteins           int      `json:"proteins" db:"proteins"`
	Fats               int      `json:"fats" db:"fats"`
	Carbohydrates      int      `json:"carbohydrates" db:"carbohydrates"`
	DietaryInformation []string `json:"dietary_information" db:"dietary_information"`
}

func (fi *FoodInsights) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}

	if len(source) == 0 {
		source = []byte("{}")
	}

	if err := json.Unmarshal(source, fi); err != nil {
		return fmt.Errorf("error unmarshalling FoodInsights JSON: %w", err)
	}

	return nil
}

func (fi FoodInsights) Value() (driver.Value, error) {
	return json.Marshal(fi)
}

type ArrayString []string

type Ingredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type Ingredients []Ingredient

func (i *Ingredients) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}

	if len(source) == 0 {
		source = []byte("[]")
	}

	if err := json.Unmarshal(source, i); err != nil {
		return fmt.Errorf("error unmarshalling Ingredients JSON: %w", err)
	}

	return nil
}

func (i Ingredients) Value() (driver.Value, error) {
	if len(i) == 0 {
		return nil, nil
	}

	return json.Marshal(i)
}

func (as *ArrayString) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		*as = make([]string, 0)
	case []byte:
		*as = strings.Split(string(src), ";")
	case string:
		*as = strings.Split(src, ";")
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
	return nil
}

func (as ArrayString) Value() (driver.Value, error) {
	if len(as) == 0 {
		return nil, nil
	}

	return strings.Join(as, ";"), nil
}

type TagSlice []Tag

func (ts *TagSlice) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	case nil:
		return json.Unmarshal([]byte("[]"), ts)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}

	if err := json.Unmarshal(source, ts); err != nil {
		return fmt.Errorf("error unmarshalling TagSlice JSON: %w", err)
	}

	return nil
}

func (s Storage) GetPostByID(uid, id int64) (*Post, error) {
	var post Post

	query := `
		SELECT p.id,
			   p.user_id,
			   p.text,
			   p.created_at,
			   p.updated_at,
			   p.hidden_at,
			   p.photo_url,
			   p.suggested_dish_name,
			   p.suggested_ingredients,
			   p.suggested_tags,
			   p.is_spam,
			   p.food_insights,
			   json_object('id', u.id, 'last_name', u.last_name, 'first_name', u.first_name, 'avatar_url', u.avatar_url, 'title',
											   u.title, 'username', u.username) AS user
		FROM posts p
				 JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`

	err := s.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Text,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.HiddenAt,
		&post.PhotoURL,
		&post.SuggestedDishName,
		&post.SuggestedIngredients,
		&post.SuggestedTags,
		&post.IsSpam,
		&post.FoodInsights,
		&post.User,
	)

	if IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s Storage) CreatePost(uid int64, post Post) (*Post, error) {
	postQuery := `
        INSERT INTO posts (user_id, photo_url, hidden_at, text)
        VALUES (?, ?, CURRENT_TIMESTAMP, ?)
    `

	res, err := s.db.Exec(postQuery, uid, post.PhotoURL, post.Text)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	post.ID = id

	return s.GetPostByID(uid, id)
}

func (s Storage) ListPosts(uid int64) ([]Post, error) {
	var posts []Post

	query := `
		SELECT p.id,
			   p.user_id,
			   p.text,
			   p.created_at,
			   p.updated_at,
			   p.hidden_at,
			   p.photo_url,
			   p.suggested_dish_name,
			   p.suggested_ingredients,
			   p.suggested_tags,
			   p.is_spam,
			   p.food_insights,
			   json_object('id', u.id, 'last_name', u.last_name, 'first_name', u.first_name, 'avatar_url', u.avatar_url,
						   'title',
						   u.title, 'username', u.username)                                                        AS user,
			   json_group_array(distinct json_object('id', t.id, 'name', t.name)) filter ( where t.id is not null) AS tags
		FROM posts p
				 JOIN users u ON p.user_id = u.id
				 LEFT JOIN post_tags pt ON p.id = pt.post_id
				 LEFT JOIN tags t ON pt.tag_id = t.id
		GROUP BY p.id
		ORDER BY p.created_at DESC
		LIMIT 100
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID,
			&p.UserID,
			&p.Text,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.HiddenAt,
			&p.PhotoURL,
			&p.SuggestedDishName,
			&p.SuggestedIngredients,
			&p.SuggestedTags,
			&p.IsSpam,
			&p.FoodInsights,
			&p.User,
			&p.Tags,
		); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

type Tag struct {
	Name string `db:"name" json:"name"`
	ID   int64  `db:"id" json:"id"`
}

func (s Storage) ListTags() ([]Tag, error) {
	var tags []Tag

	query := `
		SELECT id, name
		FROM tags
	`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s Storage) UpdatePostSuggestions(uid, postID int64, post Post) (*Post, error) {
	updateQuery := `
		UPDATE posts
		SET suggested_dish_name = ?, suggested_ingredients = ?, suggested_tags = ?, is_spam = ?
		WHERE id = ? AND user_id = ?
	`

	res, err := s.db.Exec(updateQuery, post.SuggestedDishName, post.SuggestedIngredients, post.SuggestedTags, post.IsSpam, postID, uid)
	if err != nil {
		return nil, err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return s.GetPostByID(uid, postID)
}

func (s Storage) UpdatePostFoodInsights(uid, postID int64, insights FoodInsights) (*Post, error) {
	updateQuery := `
		UPDATE posts
		SET food_insights = ?
		WHERE id = ? AND user_id = ?
	`

	_, err := s.db.Exec(updateQuery, insights, postID, uid)
	if err != nil {
		return nil, err
	}

	return s.GetPostByID(uid, postID)
}

func (s Storage) UpdatePost(uid, postID int64, post Post, tags []int) (*Post, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	updateQuery := `
        UPDATE posts
        SET text = ?, photo_url = ?, updated_at = CURRENT_TIMESTAMP, hidden_at = NULL
        WHERE id = ? AND user_id = ?
    `

	_, err = tx.Exec(updateQuery, post.Text, post.PhotoURL, postID, uid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(tags) > 0 {
		deleteQuery := `
            DELETE FROM post_tags
            WHERE post_id = ?
        `

		_, err = tx.Exec(deleteQuery, postID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tagQuery := `
            INSERT INTO post_tags (post_id, tag_id)
            VALUES (?, ?)
        `

		for _, tag := range tags {
			_, err = tx.Exec(tagQuery, postID, tag)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetPostByID(uid, postID)
}
