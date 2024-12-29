package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Meal struct {
	ID              int64         `db:"id" json:"id"`
	UserID          int64         `db:"user_id" json:"user_id"`
	Text            *string       `db:"text" json:"text"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
	HiddenAt        *time.Time    `db:"hidden_at" json:"hidden_at"`
	PhotoURL        string        `db:"photo_url" json:"photo_url"`
	DishName        *string       `json:"dish_name" db:"dish_name"`
	HealthRating    *int          `json:"health_rating" db:"health_rating"`
	AestheticRating *int          `json:"aesthetic_rating" db:"aesthetic_rating"`
	Ingredients     Ingredients   `json:"ingredients" db:"ingredients"`
	Tags            ArrayString   `json:"tags" db:"tags"`
	IsSpam          bool          `json:"is_spam" db:"is_spam"`
	FoodInsights    *FoodInsights `json:"food_insights" db:"food_insights"`
}

type FoodInsights struct {
	Calories      int `json:"calories" db:"calories"`
	Proteins      int `json:"proteins" db:"proteins"`
	Fats          int `json:"fats" db:"fats"`
	Carbohydrates int `json:"carbohydrates" db:"carbohydrates"`
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
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
	Weight   float64 `json:"weight"`
	Macros   struct {
		Proteins      float64 `json:"proteins"`
		Fats          float64 `json:"fats"`
		Carbohydrates float64 `json:"carbohydrates"`
	} `json:"macronutrients"`
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

func (s *storage) GetMealByID(id int64) (*Meal, error) {
	var meal Meal

	query := `
		SELECT m.id,
			   m.user_id,
			   m.text,
			   m.created_at,
			   m.updated_at,
			   m.hidden_at,
			   m.photo_url,
			   m.dish_name,
			   m.ingredients,
			   m.tags,
			   m.is_spam,
			   m.food_insights
		FROM meals m
		WHERE m.id = ?
	`

	err := s.db.QueryRow(query, id).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.Text,
		&meal.CreatedAt,
		&meal.UpdatedAt,
		&meal.HiddenAt,
		&meal.PhotoURL,
		&meal.DishName,
		&meal.Ingredients,
		&meal.Tags,
		&meal.IsSpam,
		&meal.FoodInsights,
	)

	if IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &meal, nil
}

func (s *storage) AddMeal(uid int64, meal Meal) (*Meal, error) {
	mealQuery := `
        INSERT INTO meals (user_id, photo_url, hidden_at, text)
        VALUES (?, ?, CURRENT_TIMESTAMP, ?)
    `

	res, err := s.db.Exec(mealQuery, uid, meal.PhotoURL, meal.Text)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	meal.ID = id

	return s.GetMealByID(id)
}

func (s *storage) ListMeals(startDate, endDate time.Time) ([]Meal, error) {
	var meals []Meal
	args := []interface{}{startDate, endDate}

	query := `
		SELECT m.id,
			   m.user_id,
			   m.text,
			   m.created_at,
			   m.updated_at,
			   m.hidden_at,
			   m.photo_url,
			   m.dish_name,
			   m.ingredients,
			   m.tags,
			   m.is_spam,
			   m.food_insights,
			   m.aesthetic_rating,
			   m.health_rating,
			   json_group_array(distinct json_object('id', t.id, 'name', t.name)) filter ( where t.id is not null) AS tags
		FROM meals m
				 JOIN users u ON m.user_id = u.id
				 LEFT JOIN meal_tags pt ON m.id = pt.meal_id
				 LEFT JOIN tags t ON pt.tag_id = t.id
		WHERE m.created_at >= ? AND m.created_at <= ?
		GROUP BY m.id
		ORDER BY m.created_at DESC
	`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m Meal
		if err := rows.Scan(&m.ID,
			&m.UserID,
			&m.Text,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.HiddenAt,
			&m.PhotoURL,
			&m.DishName,
			&m.Ingredients,
			&m.Tags,
			&m.IsSpam,
			&m.FoodInsights,
			&m.AestheticRating,
			&m.HealthRating,
			&m.Tags,
		); err != nil {
			return nil, err
		}

		meals = append(meals, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return meals, nil
}

type Tag struct {
	Name string `db:"name" json:"name"`
	ID   int64  `db:"id" json:"id"`
}

func (s *storage) ListTags() ([]Tag, error) {
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

func (s *storage) UpdateMeal(uid, mealID int64, meal Meal, tags []int) (*Meal, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	updateQuery := `
        UPDATE meals
        SET text = ?, photo_url = ?, updated_at = CURRENT_TIMESTAMP, hidden_at = NULL,
            dish_name = ?, ingredients = ?, is_spam = ?, food_insights = ?, aesthetic_rating = ?, health_rating = ?
        WHERE id = ? AND user_id = ?
    `

	_, err = tx.Exec(updateQuery, meal.Text, meal.PhotoURL, meal.DishName, meal.Ingredients, meal.IsSpam, meal.FoodInsights, meal.AestheticRating, meal.HealthRating, mealID, uid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(tags) > 0 {
		deleteQuery := `
            DELETE FROM meal_tags
            WHERE meal_id = ?
        `

		_, err = tx.Exec(deleteQuery, mealID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tagQuery := `
            INSERT INTO meal_tags (meal_id, tag_id)
            VALUES (?, ?)
        `

		for _, tag := range tags {
			_, err = tx.Exec(tagQuery, mealID, tag)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetMealByID(mealID)
}
