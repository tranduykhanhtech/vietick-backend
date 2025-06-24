package repository

import (
	"fmt"

	"vietick-backend/internal/model"
	"vietick-backend/internal/utils"

	"gorm.io/gorm"
)

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(verification *model.IdentityVerification) error {
	if err := r.db.Create(verification).Error; err != nil {
		return fmt.Errorf("failed to create identity verification: %w", err)
	}
	return nil
}

func (r *VerificationRepository) GetByID(verificationID string) (*model.IdentityVerification, error) {
	verification := &model.IdentityVerification{}
	if err := r.db.First(verification, verificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("verification not found")
		}
		return nil, fmt.Errorf("failed to get verification: %w", err)
	}
	return verification, nil
}

func (r *VerificationRepository) GetByUserID(userID string) (*model.IdentityVerification, error) {
	verification := &model.IdentityVerification{}
	if err := r.db.Where("user_id = ?", userID).Order("submitted_at DESC").First(verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("verification not found")
		}
		return nil, fmt.Errorf("failed to get verification: %w", err)
	}
	return verification, nil
}

func (r *VerificationRepository) GetPending(pagination utils.PaginationResult) ([]model.IdentityVerification, int64, error) {
	var verifications []model.IdentityVerification
	var totalCount int64
	r.db.Model(&model.IdentityVerification{}).Where("status = ?", "pending").Count(&totalCount)
	if err := r.db.Where("status = ?", "pending").Order("submitted_at ASC").Limit(pagination.Limit).Offset(pagination.Offset).Find(&verifications).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get pending verifications: %w", err)
	}
	return verifications, totalCount, nil
}

func (r *VerificationRepository) Review(verificationID string, status model.IdentityVerificationStatus, adminNotes *string, reviewedBy string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.IdentityVerification{}).Where("id = ?", verificationID).Updates(map[string]interface{}{
			"status":      status,
			"admin_notes": adminNotes,
			"reviewed_at": gorm.Expr("NOW()"),
			"reviewed_by": reviewedBy,
		}).Error; err != nil {
			return fmt.Errorf("failed to update verification: %w", err)
		}
		if status == model.IdentityVerificationApproved {
			var verification model.IdentityVerification
			if err := tx.First(&verification, verificationID).Error; err != nil {
				return fmt.Errorf("failed to get verification: %w", err)
			}
			if err := tx.Model(&model.User{}).Where("id = ?", verification.UserID).Updates(map[string]interface{}{
				"is_verified":                  true,
				"identity_verification_status": "approved",
				"updated_at":                   gorm.Expr("NOW()"),
			}).Error; err != nil {
				return fmt.Errorf("failed to update user verification status: %w", err)
			}
		}
		return nil
	})
}

func (r *VerificationRepository) GetAll(status *model.IdentityVerificationStatus, pagination utils.PaginationResult) ([]model.IdentityVerification, int64, error) {
	var verifications []model.IdentityVerification
	var totalCount int64
	dbQuery := r.db.Model(&model.IdentityVerification{})
	if status != nil {
		dbQuery = dbQuery.Where("status = ?", *status)
	}
	dbQuery.Count(&totalCount)
	if err := dbQuery.Order("submitted_at DESC").Limit(pagination.Limit).Offset(pagination.Offset).Find(&verifications).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get verifications: %w", err)
	}
	return verifications, totalCount, nil
}

func (r *VerificationRepository) HasPendingVerification(userID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.IdentityVerification{}).
		Where("user_id = ? AND status = ?", userID, "pending").
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check pending verification: %w", err)
	}
	return count > 0, nil
}

func (r *VerificationRepository) Delete(verificationID string) error {
	result := r.db.Where("id = ?", verificationID).Delete(&model.IdentityVerification{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete verification: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("verification not found")
	}
	return nil
}
