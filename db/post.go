package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type PostReactions struct {
	Frown int `json:"frown" db:"frown"`
	Smile int `json:"smile" db:"smile"`
	Meh   int `json:"meh" db:"meh"`
}

type UserReaction struct {
	Type       string `db:"type" json:"type"`
	HasReacted bool   `db:"has_reacted" json:"has_reacted"`
}

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

type Location struct {
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Address   string    `json:"address" db:"address"`
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Post struct {
	ID                   int64         `db:"id" json:"id"`
	UserID               int64         `db:"user_id" json:"user_id"`
	Text                 *string       `db:"text" json:"text"`
	LocationID           *int64        `db:"location_id" json:"location_id"`
	CreatedAt            time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time     `db:"updated_at" json:"updated_at"`
	HiddenAt             *time.Time    `db:"hidden_at" json:"hidden_at"`
	PhotoURL             string        `db:"photo_url" json:"photo_url"`
	User                 UserProfile   `json:"user" db:"user"`
	DishName             *string       `json:"dish_name" db:"dish_name"`
	Ingredients          ArrayString   `json:"ingredients" db:"ingredients"`
	SuggestedDishName    *string       `json:"suggested_dish_name" db:"suggested_dish_name"`
	SuggestedIngredients ArrayString   `json:"suggested_ingredients" db:"suggested_ingredients"`
	SuggestedTags        ArrayString   `json:"suggested_tags" db:"suggested_tags"`
	IsSpam               bool          `json:"is_spam" db:"is_spam"`
	Location             *Location     `json:"location" db:"location"`
	Reactions            PostReactions `json:"reactions" db:"-"`
	UserReaction         UserReaction  `json:"user_reaction" db:"-"`
	Tags                 TagSlice      `json:"tags" db:"tags"`
}

type ArrayString []string

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
			   p.ingredients,
			   p.dish_name,
			   p.suggested_dish_name,
			   p.suggested_ingredients,
			   p.suggested_tags,
			   p.is_spam,
			   p.location_id,
			   json_object('id', u.id, 'last_name', u.last_name, 'first_name', u.first_name, 'avatar_url', u.avatar_url, 'title',
											   u.title, 'username', u.username) AS user,
			   COALESCE(SUM(CASE WHEN r.type = 'smile' THEN 1 ELSE 0 END), 0) AS smile,
			   COALESCE(SUM(CASE WHEN r.type = 'frown' THEN 1 ELSE 0 END), 0) AS frown,
			   COALESCE(SUM(CASE WHEN r.type = 'meh' THEN 1 ELSE 0 END), 0) AS meh,
			   COALESCE(r2.type, '') AS type,
			   json_object('id', l.id, 'latitude', l.latitude, 'longitude', l.longitude, 'address', l.address, 'user_id', l.user_id) AS location
		FROM posts p
				 JOIN users u ON p.user_id = u.id
				 LEFT JOIN reactions r ON p.id = r.post_id
				 LEFT JOIN reactions r2 ON p.id = r2.post_id AND r2.user_id = ?
				 LEFT JOIN locations l ON p.location_id = l.id
		WHERE p.id = ?
	`

	var locationJSON []byte

	err := s.db.QueryRow(query, uid, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Text,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.HiddenAt,
		&post.PhotoURL,
		&post.Ingredients,
		&post.DishName,
		&post.SuggestedDishName,
		&post.SuggestedIngredients,
		&post.SuggestedTags,
		&post.IsSpam,
		&post.LocationID,
		&post.User,
		&post.Reactions.Smile,
		&post.Reactions.Frown,
		&post.Reactions.Meh,
		&post.UserReaction.Type,
		&locationJSON,
	)

	if IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if post.LocationID != nil {
		var location Location
		if err := json.Unmarshal(locationJSON, &location); err != nil {
			return nil, err
		}
		post.Location = &location
	}

	return &post, nil
}

func (s Storage) CreatePost(uid int64, post Post) (*Post, error) {
	postQuery := `
        INSERT INTO posts (user_id, photo_url, hidden_at)
        VALUES (?, ?, CURRENT_TIMESTAMP)
    `

	res, err := s.db.Exec(postQuery, uid, post.PhotoURL)
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
			   p.ingredients,
			   p.dish_name,
			   p.suggested_dish_name,
			   p.suggested_ingredients,
			   p.suggested_tags,
			   p.is_spam,
			   p.location_id,
			   json_object('id', u.id, 'last_name', u.last_name, 'first_name', u.first_name, 'avatar_url', u.avatar_url,
						   'title',
						   u.title, 'username', u.username)                                                        AS user,
			   json_group_array(distinct json_object('id', t.id, 'name', t.name)) filter ( where t.id is not null) AS tags,
			   json_object('id', l.id, 'latitude', l.latitude, 'longitude', l.longitude, 'address', l.address, 'user_id', l.user_id) AS location
		FROM posts p
				 JOIN users u ON p.user_id = u.id
				 LEFT JOIN post_tags pt ON p.id = pt.post_id
				 LEFT JOIN tags t ON pt.tag_id = t.id
				 LEFT JOIN locations l ON p.location_id = l.id
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
		var locationJSON []byte
		if err := rows.Scan(&p.ID,
			&p.UserID,
			&p.Text,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.HiddenAt,
			&p.PhotoURL,
			&p.Ingredients,
			&p.DishName,
			&p.SuggestedDishName,
			&p.SuggestedIngredients,
			&p.SuggestedTags,
			&p.IsSpam,
			&p.LocationID,
			&p.User,
			&p.Tags,
			&locationJSON,
		); err != nil {
			return nil, err
		}

		if p.LocationID != nil {
			var location Location
			if err := json.Unmarshal(locationJSON, &location); err != nil {
				return nil, err
			}
			p.Location = &location
		}

		reactionQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'smile' THEN 1 ELSE 0 END), 0) AS smile,
			COALESCE(SUM(CASE WHEN type = 'frown' THEN 1 ELSE 0 END), 0) AS frown,
			COALESCE(SUM(CASE WHEN type = 'meh' THEN 1 ELSE 0 END), 0) AS meh
		FROM reactions
		WHERE post_id = ?
		`

		var reactions PostReactions
		err = s.db.QueryRow(reactionQuery, p.ID).Scan(&reactions.Smile, &reactions.Frown, &reactions.Meh)
		if err != nil {
			return nil, err
		}

		p.Reactions = reactions

		userReactionQuery := `SELECT type FROM reactions WHERE post_id = ? AND user_id = ?`
		err = s.db.QueryRow(userReactionQuery, p.ID, uid).Scan(&p.UserReaction.Type)
		if IsNoRowsError(err) {
			p.UserReaction.HasReacted = false
		} else if err != nil {
			return nil, err
		} else {
			p.UserReaction.HasReacted = true
		}

		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

type ReactionType string

const (
	ReactionSmile ReactionType = "smile"
	ReactionMeh   ReactionType = "meh"
	ReactionFrown ReactionType = "frown"
)

func (s Storage) UpdatePostReaction(uid, postID int64, reactionType ReactionType) error {
	query := `
		INSERT INTO reactions (user_id, post_id, type)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, post_id) DO UPDATE SET type = ?, created_at = CURRENT_TIMESTAMP
	`

	_, err := s.db.Exec(query, uid, postID, reactionType, reactionType)
	return err
}

func (s Storage) DeletePostReaction(uid, postID int64) error {
	query := `
		DELETE FROM reactions
		WHERE user_id = ? AND post_id = ?
	`

	res, err := s.db.Exec(query, uid, postID)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
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

func (s Storage) UpdatePost(uid, postID int64, post Post, tags []int) (*Post, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var locationID *int64

	if post.Location != nil {
		locationQuery := `
            INSERT INTO locations (latitude, longitude, address, user_id)
            VALUES (?, ?, ?, ?)
        `

		res, err := tx.Exec(locationQuery, post.Location.Latitude, post.Location.Longitude, post.Location.Address, uid)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		locationID = &id
	}

	updateQuery := `
        UPDATE posts
        SET text = ?, photo_url = ?, location_id = ?, updated_at = CURRENT_TIMESTAMP, hidden_at = NULL
        WHERE id = ? AND user_id = ?
    `

	_, err = tx.Exec(updateQuery, post.Text, post.PhotoURL, post.SuggestedDishName, post.SuggestedIngredients, post.SuggestedTags, post.IsSpam, locationID, postID, uid)
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
