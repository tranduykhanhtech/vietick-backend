package repository

import (
	"fmt"

	"vietick-backend/internal/model"
	"vietick-backend/internal/utils"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
		if err := tx.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
			return fmt.Errorf("failed to update comment count: %w", err)
		}
		return nil
	})
}

func (r *CommentRepository) GetByID(commentID string, userID *string) (*model.Comment, error) {
	comment := &model.Comment{}
	if err := r.db.First(comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	user := &model.UserProfile{}
	r.db.Model(&model.User{}).Select("id, username, full_name, bio, avatar_url, is_verified, created_at").Where("id = ?", comment.UserID).Scan(user)
	comment.User = user
	if userID != nil {
		var count int64
		r.db.Model(&model.Comment{}).Where("id = ? AND user_id = ?", commentID, *userID).Count(&count)
		comment.IsLiked = count > 0
	}
	return comment, nil
}

func (r *CommentRepository) GetPostComments(postID string, userID *string, pagination utils.PaginationResult) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var totalCount int64
	r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&totalCount)
	if err := r.db.Where("post_id = ?", postID).Order("created_at ASC").Limit(pagination.Limit).Offset(pagination.Offset).Find(&comments).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}
	return comments, totalCount, nil
}

func (r *CommentRepository) Update(comment *model.Comment) error {
	if err := r.db.Save(comment).Error; err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

func (r *CommentRepository) Delete(commentID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		if err := tx.Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("comment not found or you're not the owner")
			}
			return fmt.Errorf("failed to get comment: %w", err)
		}
		if err := tx.Delete(&comment).Error; err != nil {
			return fmt.Errorf("failed to delete comment: %w", err)
		}
		if err := tx.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
			return fmt.Errorf("failed to update comment count: %w", err)
		}
		return nil
	})
}

func (r *CommentRepository) LikeComment(commentID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert like
		res := tx.Exec("INSERT IGNORE INTO comment_likes (comment_id, user_id) VALUES (?, ?)", commentID, userID)
		if res.Error != nil {
			return fmt.Errorf("failed to like comment: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("comment already liked")
		}
		// Update like count
		if err := tx.Model(&model.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return fmt.Errorf("failed to update like count: %w", err)
		}
		return nil
	})
}

func (r *CommentRepository) UnlikeComment(commentID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete like
		res := tx.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if res.Error != nil {
			return fmt.Errorf("failed to unlike comment: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("comment not liked")
		}
		// Update like count
		if err := tx.Model(&model.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
			return fmt.Errorf("failed to update like count: %w", err)
		}
		return nil
	})
}

func (r *CommentRepository) IsCommentLikedByUser(commentID, userID string) (bool, error) {
	var count int64
	err := r.db.Table("comment_likes").
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if comment is liked: %w", err)
	}
	return count > 0, nil
}
