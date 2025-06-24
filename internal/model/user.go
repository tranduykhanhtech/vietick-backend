package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID                         string                     `json:"id" db:"id"`
	Username                   string                     `json:"username" db:"username"`
	Email                      string                     `json:"email" db:"email"`
	PasswordHash               string                     `json:"-" db:"password_hash"`
	FullName                   string                     `json:"full_name" db:"full_name"`
	Bio                        *string                    `json:"bio" db:"bio"`
	AvatarURL                  *string                    `json:"avatar_url" db:"avatar_url"`
	IsVerified                 bool                       `json:"is_verified" db:"is_verified"`
	IsEmailVerified            bool                       `json:"is_email_verified" db:"is_email_verified"`
	EmailVerificationToken     *string                    `json:"-" db:"email_verification_token"`
	EmailVerificationExpiresAt *time.Time                 `json:"-" db:"email_verification_expires_at"`
	IdentityVerificationStatus IdentityVerificationStatus `json:"identity_verification_status" db:"identity_verification_status"`
	IdentityDocuments          *IdentityDocuments         `json:"identity_documents" db:"identity_documents"`
	CreatedAt                  time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                  `json:"updated_at" db:"updated_at"`
}

type IdentityVerificationStatus string

const (
	IdentityVerificationNone     IdentityVerificationStatus = "none"
	IdentityVerificationPending  IdentityVerificationStatus = "pending"
	IdentityVerificationApproved IdentityVerificationStatus = "approved"
	IdentityVerificationRejected IdentityVerificationStatus = "rejected"
)

type IdentityDocuments struct {
	FrontImageURL  string `json:"front_image_url"`
	BackImageURL   string `json:"back_image_url"`
	SelfieImageURL string `json:"selfie_image_url"`
}

// Implement sql.Scanner interface for JSON fields
func (id *IdentityDocuments) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into IdentityDocuments", value)
	}

	return json.Unmarshal(bytes, id)
}

// Implement driver.Valuer interface for JSON fields
func (id IdentityDocuments) Value() (driver.Value, error) {
	return json.Marshal(id)
}

// UserProfile represents public user information
type UserProfile struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Bio            *string   `json:"bio"`
	AvatarURL      *string   `json:"avatar_url"`
	IsVerified     bool      `json:"is_verified"`
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
	PostsCount     int       `json:"posts_count"`
	IsFollowing    bool      `json:"is_following,omitempty"`
	IsFollowedBy   bool      `json:"is_followed_by,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=1,max=100"`
}

// LoginRequest represents user login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents profile update data
type UpdateProfileRequest struct {
	FullName  *string `json:"full_name,omitempty" binding:"omitempty,min=1,max=100"`
	Bio       *string `json:"bio,omitempty" binding:"omitempty,max=500"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
