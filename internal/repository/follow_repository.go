package repository

import (
	"fmt"

	"vietick-backend/internal/model"
	"vietick-backend/internal/utils"

	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) Follow(followerID, followingID string) error {
	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}
	// Kiểm tra đã follow chưa
	var count int64
	r.db.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count)
	if count > 0 {
		return fmt.Errorf("already following this user")
	}
	follow := &model.Follow{FollowerID: followerID, FollowingID: followingID}
	if err := r.db.Create(follow).Error; err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}
	return nil
}

func (r *FollowRepository) Unfollow(followerID, followingID string) error {
	result := r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&model.Follow{})
	if result.Error != nil {
		return fmt.Errorf("failed to unfollow user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("not following this user")
	}
	return nil
}

func (r *FollowRepository) IsFollowing(followerID, followingID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}
	return count > 0, nil
}

func (r *FollowRepository) GetFollowers(userID string, viewerID *string, pagination utils.PaginationResult) ([]model.UserProfile, int64, error) {
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
	query += `
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.following_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`

	countQuery := `SELECT COUNT(*) FROM follows WHERE following_id = ?`
	var totalCount int64
	if err := r.db.Raw(countQuery, userID).Scan(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	var users []model.UserProfile
	var args []interface{}
	if viewerID != nil {
		args = append(args, *viewerID, *viewerID, userID, pagination.Limit, pagination.Offset)
	} else {
		args = append(args, userID, pagination.Limit, pagination.Offset)
	}
	if err := r.db.Raw(query, args...).Scan(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get followers: %w", err)
	}
	return users, totalCount, nil
}

func (r *FollowRepository) GetFollowing(userID string, viewerID *string, pagination utils.PaginationResult) ([]model.UserProfile, int64, error) {
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
	query += `
		FROM users u
		JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`

	countQuery := `SELECT COUNT(*) FROM follows WHERE follower_id = ?`
	var totalCount int64
	if err := r.db.Raw(countQuery, userID).Scan(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count following: %w", err)
	}

	var users []model.UserProfile
	var args []interface{}
	if viewerID != nil {
		args = append(args, *viewerID, *viewerID, userID, pagination.Limit, pagination.Offset)
	} else {
		args = append(args, userID, pagination.Limit, pagination.Offset)
	}
	if err := r.db.Raw(query, args...).Scan(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get following: %w", err)
	}
	return users, totalCount, nil
}

func (r *FollowRepository) GetFollowCounts(userID string) (followersCount, followingCount int64, err error) {
	err = r.db.Model(&model.Follow{}).Where("following_id = ?", userID).Count(&followersCount).Error
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get followers count: %w", err)
	}
	err = r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&followingCount).Error
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get following count: %w", err)
	}
	return followersCount, followingCount, nil
}

func (r *FollowRepository) GetMutualFollows(userID1, userID2 string, limit int) ([]model.UserProfile, error) {
	query := `
		SELECT u.id, u.username, u.full_name, u.bio, u.avatar_url, u.is_verified, u.created_at,
		       (SELECT COUNT(*) FROM follows WHERE following_id = u.id) as followers_count,
		       (SELECT COUNT(*) FROM follows WHERE follower_id = u.id) as following_count,
		       (SELECT COUNT(*) FROM posts WHERE user_id = u.id) as posts_count
		FROM users u
		WHERE u.id IN (
		    SELECT f1.following_id 
		    FROM follows f1
		    JOIN follows f2 ON f1.following_id = f2.following_id
		    WHERE f1.follower_id = ? AND f2.follower_id = ?
		)
		ORDER BY u.created_at DESC
		LIMIT ?
	`

	var users []model.UserProfile
	if err := r.db.Raw(query, userID1, userID2, limit).Scan(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get mutual follows: %w", err)
	}
	return users, nil
}
