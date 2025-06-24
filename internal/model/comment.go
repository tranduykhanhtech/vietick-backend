package model

import "time"

type Comment struct {
	ID        string     `json:"id" db:"id" gorm:"type:char(36)"`
	PostID    string     `json:"post_id" db:"post_id" gorm:"type:char(36)"`
	UserID    string     `json:"user_id" db:"user_id" gorm:"type:char(36)"`
	Content   string     `json:"content" db:"content"`
	LikeCount int        `json:"like_count" db:"like_count"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`

	// Additional fields for API responses
	User    *UserProfile `json:"user,omitempty"`
	IsLiked bool         `json:"is_liked,omitempty"`
}

type CommentLike struct {
	ID        string    `json:"id" db:"id"`
	CommentID string    `json:"comment_id" db:"comment_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Request models
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

type CommentsResponse struct {
	Comments   []Comment `json:"comments"`
	TotalCount int64     `json:"total_count"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	HasMore    bool      `json:"has_more"`
} 