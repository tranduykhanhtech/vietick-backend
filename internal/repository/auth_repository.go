package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"vietick-backend/internal/model"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateRefreshToken(refreshToken *model.RefreshToken) error {
	if err := r.db.Create(refreshToken).Error; err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetRefreshToken(tokenHash string) (*model.RefreshToken, error) {
	refreshToken := &model.RefreshToken{}
	if err := r.db.Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).First(refreshToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return refreshToken, nil
}

func (r *AuthRepository) DeleteRefreshToken(tokenHash string) error {
	if err := r.db.Where("token_hash = ?", tokenHash).Delete(&model.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

func (r *AuthRepository) DeleteUserRefreshTokens(userID string) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete user refresh tokens: %w", err)
	}
	return nil
}

func (r *AuthRepository) CleanupExpiredTokens() error {
	if err := r.db.Where("expires_at <= ?", time.Now()).Delete(&model.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetUserRefreshTokenCount(userID string) (int, error) {
	var count int64
	if err := r.db.Model(&model.RefreshToken{}).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get user refresh token count: %w", err)
	}
	return int(count), nil
}

func (r *AuthRepository) DeleteOldestUserRefreshToken(userID string) error {
	var oldest model.RefreshToken
	if err := r.db.Where("user_id = ?", userID).Order("created_at ASC").First(&oldest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // Không có token nào để xóa
		}
		return fmt.Errorf("failed to find oldest refresh token: %w", err)
	}
	if err := r.db.Delete(&oldest).Error; err != nil {
		return fmt.Errorf("failed to delete oldest refresh token: %w", err)
	}
	return nil
}
