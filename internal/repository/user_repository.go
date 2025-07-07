package repository

import (
	"fmt"
	"time"

	"vietick-backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("username = ?", username).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmailVerificationToken(token string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("email_verification_token = ? AND email_verification_expires_at > ?", token, time.Now()).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired verification token")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) VerifyEmail(userID string) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_email_verified":             true,
		"email_verification_token":      nil,
		"email_verification_expires_at": nil,
		"updated_at":                    time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdatePassword(userID string, passwordHash string) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_hash": passwordHash,
		"updated_at":    time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

func (r *UserRepository) GetProfile(userID string, viewerID *string) (*model.UserProfile, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.bio, u.avatar_url, u.is_verified, u.created_at,
		       (SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers_count,
		       (SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following_count,
		       (SELECT COUNT(*) FROM posts WHERE user_id = u.id) as posts_count
	`
	if viewerID != nil {
		query += `,
		       EXISTS(SELECT 1 FROM follows WHERE follower_id = ? AND following_id = u.id) as is_following,
		       EXISTS(SELECT 1 FROM follows WHERE follower_id = u.id AND following_id = ?) as is_followed_by
		`
	}
	query += ` FROM users u WHERE u.id = ?`

	profile := &model.UserProfile{}
	var args []interface{}
	if viewerID != nil {
		args = append(args, *viewerID, *viewerID, userID)
	} else {
		args = append(args, userID)
	}
	if err := r.db.Raw(query, args...).Scan(profile).Error; err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return profile, nil
}

func (r *UserRepository) SearchUsers(query string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var totalCount int64
	q := "%" + query + "%"
	db := r.db.Model(&model.User{}).
		Where("username LIKE ? OR full_name LIKE ? OR email LIKE ?", q, q, q)
	db.Count(&totalCount)
	err := db.Order("created_at DESC").
		Limit(pageSize).
		Offset((page-1)*pageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, totalCount, nil
}

func (r *UserRepository) UpdateEmailVerificationToken(userID string, token string, expiresAt time.Time) error {
	if err := r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email_verification_token":      token,
			"email_verification_expires_at": expiresAt,
			"updated_at":                    time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update email verification token: %w", err)
	}
	return nil
}
