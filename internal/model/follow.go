package model

import "time"

type Follow struct {
	ID          string    `json:"id" db:"id"`
	FollowerID  string    `json:"follower_id" db:"follower_id"`
	FollowingID string    `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type FollowersResponse struct {
	Users      []UserProfile `json:"users"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	HasMore    bool          `json:"has_more"`
}

type FollowingResponse struct {
	Users      []UserProfile `json:"users"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	HasMore    bool          `json:"has_more"`
}
