package repository

import (
	"fmt"

	"vietick-backend/internal/model"
	"vietick-backend/internal/utils"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *PostRepository) GetByID(postID string, userID *string) (*model.Post, error) {
	post := &model.Post{}
	if err := r.db.First(post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	// Lấy thông tin user
	user := &model.UserProfile{}
	r.db.Model(&model.User{}).Select("id, username, full_name, bio, avatar_url, is_verified, created_at").Where("id = ?", post.UserID).Scan(user)
	post.User = user
	// Nếu có userID, kiểm tra like
	if userID != nil {
		var count int64
		r.db.Model(&model.PostLike{}).Where("post_id = ? AND user_id = ?", postID, *userID).Count(&count)
		post.IsLiked = count > 0
	}
	return post, nil
}

func (r *PostRepository) Update(post *model.Post) error {
	if err := r.db.Save(post).Error; err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (r *PostRepository) Delete(postID, userID string) error {
	if err := r.db.Where("id = ? AND user_id = ?", postID, userID).Delete(&model.Post{}).Error; err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (r *PostRepository) GetFeed(userID string, pagination utils.PaginationResult) ([]model.Post, int64, error) {
	var posts []model.Post
	// Lấy danh sách user_id mà user này theo dõi + chính user đó
	var followingIDs []string
	followingIDs = append(followingIDs, userID)
	var ids []string
	r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Pluck("following_id", &ids)
	followingIDs = append(followingIDs, ids...)
	// Đếm tổng số post
	var totalCount int64
	r.db.Model(&model.Post{}).Where("user_id IN ?", followingIDs).Count(&totalCount)
	// Lấy post
	if err := r.db.Where("user_id IN ?", followingIDs).Order("created_at DESC").Limit(pagination.Limit).Offset(pagination.Offset).Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get feed: %w", err)
	}
	return posts, totalCount, nil
}

func (r *PostRepository) GetUserPosts(userID string, viewerID *string, pagination utils.PaginationResult) ([]model.Post, int64, error) {
	var posts []model.Post
	var totalCount int64
	r.db.Model(&model.Post{}).Where("user_id = ?", userID).Count(&totalCount)
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(pagination.Limit).Offset(pagination.Offset).Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get user posts: %w", err)
	}
	return posts, totalCount, nil
}

func (r *PostRepository) LikePost(postID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert like
		res := tx.Exec("INSERT IGNORE INTO post_likes (post_id, user_id) VALUES (?, ?)", postID, userID)
		if res.Error != nil {
			return fmt.Errorf("failed to like post: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("post already liked")
		}
		// Update like count
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return fmt.Errorf("failed to update like count: %w", err)
		}
		return nil
	})
}

func (r *PostRepository) UnlikePost(postID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete like
		res := tx.Exec("DELETE FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)
		if res.Error != nil {
			return fmt.Errorf("failed to unlike post: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("post not liked")
		}
		// Update like count
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
			return fmt.Errorf("failed to update like count: %w", err)
		}
		return nil
	})
}

func (r *PostRepository) IsPostLikedByUser(postID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.PostLike{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if post is liked: %w", err)
	}
	return count > 0, nil
}
