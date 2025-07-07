package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Post struct {
	ID           string     `json:"id" db:"id"`
	UserID       string     `json:"user_id" db:"user_id"`
	Content      string     `json:"content" db:"content"`
	ImageURLs    ImageURLs  `json:"image_urls" db:"image_urls" gorm:"type:json"` // Thêm tag này
	LikeCount    int        `json:"like_count" db:"like_count"`
	CommentCount int        `json:"comment_count" db:"comment_count"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	// Additional fields for API responses
	User    *UserProfile `json:"user,omitempty"`
	IsLiked bool         `json:"is_liked,omitempty"`
}

type ImageURLs []string

// Implement sql.Scanner interface for JSON fields
func (iu *ImageURLs) Scan(value interface{}) error {
	if value == nil {
		*iu = ImageURLs{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into ImageURLs", value)
	}

	return json.Unmarshal(bytes, iu)
}

// Implement driver.Valuer interface for JSON fields
func (iu ImageURLs) Value() (driver.Value, error) {
	if len(iu) == 0 {
		return nil, nil
	}
	return json.Marshal(iu)
}

type PostLike struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"post_id" db:"post_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Request models
type CreatePostRequest struct {
	Content   string   `json:"content" binding:"required,min=1,max=5000"`
	ImageURLs []string `json:"image_urls,omitempty"`
}

type UpdatePostRequest struct {
	Content   string   `json:"content" binding:"required,min=1,max=5000"`
	ImageURLs []string `json:"image_urls,omitempty"`
}

type PostsResponse struct {
	Posts      []Post `json:"posts"`
	TotalCount int64  `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	HasMore    bool   `json:"has_more"`
}

type Hashtag struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PostHashtag struct {
	PostID    string    `json:"post_id" db:"post_id"`
	HashtagID string    `json:"hashtag_id" db:"hashtag_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
