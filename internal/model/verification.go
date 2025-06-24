package model

import "time"

type IdentityVerification struct {
	ID             string                     `json:"id" db:"id"`
	UserID         string                     `json:"user_id" db:"user_id"`
	FullName       string                     `json:"full_name" db:"full_name"`
	IDNumber       string                     `json:"id_number" db:"id_number"`
	IDType         IdentityDocumentType       `json:"id_type" db:"id_type"`
	FrontImageURL  string                     `json:"front_image_url" db:"front_image_url"`
	BackImageURL   *string                    `json:"back_image_url" db:"back_image_url"`
	SelfieImageURL string                     `json:"selfie_image_url" db:"selfie_image_url"`
	Status         IdentityVerificationStatus `json:"status" db:"status"`
	AdminNotes     *string                    `json:"admin_notes" db:"admin_notes"`
	SubmittedAt    time.Time                  `json:"submitted_at" db:"submitted_at"`
	ReviewedAt     *time.Time                 `json:"reviewed_at" db:"reviewed_at"`
	ReviewedBy     *string                    `json:"reviewed_by" db:"reviewed_by"`

	// Additional fields for API responses
	User     *UserProfile `json:"user,omitempty" gorm:"-"`
	Reviewer *UserProfile `json:"reviewer,omitempty" gorm:"-"`
}

type IdentityDocumentType string

const (
	IdentityDocumentNationalID    IdentityDocumentType = "national_id"
	IdentityDocumentPassport      IdentityDocumentType = "passport"
	IdentityDocumentDriverLicense IdentityDocumentType = "driver_license"
)

type SubmitIdentityVerificationRequest struct {
	FullName       string               `json:"full_name" binding:"required,min=1,max=100"`
	IDNumber       string               `json:"id_number" binding:"required,min=1,max=50"`
	IDType         IdentityDocumentType `json:"id_type" binding:"required,oneof=national_id passport driver_license"`
	FrontImageURL  string               `json:"front_image_url" binding:"required"`
	BackImageURL   *string              `json:"back_image_url,omitempty"`
	SelfieImageURL string               `json:"selfie_image_url" binding:"required"`
}

type ReviewIdentityVerificationRequest struct {
	Status     IdentityVerificationStatus `json:"status" binding:"required,oneof=approved rejected"`
	AdminNotes *string                    `json:"admin_notes,omitempty"`
}

type IdentityVerificationsResponse struct {
	Verifications []IdentityVerification `json:"verifications"`
	TotalCount    int64                  `json:"total_count"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
	HasMore       bool                   `json:"has_more"`
}
