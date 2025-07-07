package repository

import (
	"fmt"
	"time"

	"vietick-backend/internal/model"
	"vietick-backend/internal/utils"

	"gorm.io/gorm"
	"github.com/google/uuid"
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

// Hashtag repository methods
func (r *PostRepository) FindOrCreateHashtag(name string) (*model.Hashtag, error) {
	hashtag := &model.Hashtag{}
	err := r.db.Where("name = ?", name).First(hashtag).Error
	if err == nil {
		return hashtag, nil
	}
	if err.Error() == "record not found" || err.Error() == "gorm: record not found" {
		hashtag.ID = uuid.New().String()
		hashtag.Name = name
		hashtag.CreatedAt = time.Now()
		if err := r.db.Create(hashtag).Error; err != nil {
			return nil, err
		}
		return hashtag, nil
	}
	return nil, err
}

func (r *PostRepository) AddHashtagsToPost(postID string, hashtagNames []string) error {
	for _, name := range hashtagNames {
		hashtag, err := r.FindOrCreateHashtag(name)
		if err != nil {
			return err
		}
		postHashtag := &model.PostHashtag{
			PostID:    postID,
			HashtagID: hashtag.ID,
			CreatedAt: time.Now(),
		}
		// Sử dụng INSERT IGNORE để tránh duplicate
		err = r.db.Exec("INSERT IGNORE INTO post_hashtags (post_id, hashtag_id, created_at) VALUES (?, ?, ?)", postID, hashtag.ID, postHashtag.CreatedAt).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepository) GetHashtagsByPost(postID string) ([]model.Hashtag, error) {
	var hashtags []model.Hashtag
	err := r.db.Raw(`SELECT h.* FROM hashtags h JOIN post_hashtags ph ON h.id = ph.hashtag_id WHERE ph.post_id = ?`, postID).Scan(&hashtags).Error
	if err != nil {
		return nil, err
	}
	return hashtags, nil
}

func (r *PostRepository) GetPostsByHashtag(hashtagName string, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Raw(`SELECT p.* FROM posts p JOIN post_hashtags ph ON p.id = ph.post_id JOIN hashtags h ON ph.hashtag_id = h.id WHERE h.name = ? ORDER BY p.created_at DESC LIMIT ? OFFSET ?`, hashtagName, limit, offset).Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// SearchPosts tìm kiếm post theo content, hashtag, username, full_name
func (r *PostRepository) SearchPosts(query string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var totalCount int64
	q := "%" + query + "%"

	db := r.db.Table("posts p").
		Select("DISTINCT p.*").
		Joins("LEFT JOIN users u ON p.user_id = u.id").
		Joins("LEFT JOIN post_hashtags ph ON p.id = ph.post_id").
		Joins("LEFT JOIN hashtags h ON ph.hashtag_id = h.id").
		Where("p.content LIKE ? OR h.name LIKE ? OR u.username LIKE ? OR u.full_name LIKE ?", q, q, q, q)

	db.Count(&totalCount)

	err := db.Order("p.created_at DESC").
		Limit(pageSize).
		Offset((page-1)*pageSize).
		Scan(&posts).Error
	if err != nil {
		return nil, 0, err
	}
	return posts, totalCount, nil
}

// SearchHashtags tìm kiếm hashtag theo tên
func (r *PostRepository) SearchHashtags(query string, page, pageSize int) ([]model.Hashtag, int64, error) {
	var hashtags []model.Hashtag
	var totalCount int64
	q := "%" + query + "%"
	db := r.db.Model(&model.Hashtag{}).
		Where("name LIKE ?", q)
	db.Count(&totalCount)
	err := db.Order("created_at DESC").
		Limit(pageSize).
		Offset((page-1)*pageSize).
		Find(&hashtags).Error
	if err != nil {
		return nil, 0, err
	}
	return hashtags, totalCount, nil
}

// SearchPostsByContent tìm kiếm post chỉ theo nội dung content
func (r *PostRepository) SearchPostsByContent(query string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var totalCount int64
	q := "%" + query + "%"
	db := r.db.Model(&model.Post{}).
		Where("content LIKE ?", q)
	db.Count(&totalCount)
	err := db.Order("created_at DESC").
		Limit(pageSize).
		Offset((page-1)*pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}
	return posts, totalCount, nil
}

// Xóa tất cả hashtag liên kết với post
func (r *PostRepository) ClearPostHashtags(postID string) error {
	return r.db.Exec("DELETE FROM post_hashtags WHERE post_id = ?", postID).Error
}
