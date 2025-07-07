package service

import (
	"fmt"
	"regexp"
	"strings"

	"vietick-backend/internal/model"
	"vietick-backend/internal/repository"
	"vietick-backend/internal/utils"
	"github.com/google/uuid"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func extractHashtags(content string) []string {
	hashtagRegex := regexp.MustCompile(`#([\p{L}0-9_]+)`)
	matches := hashtagRegex.FindAllStringSubmatch(content, -1)
	unique := map[string]struct{}{}
	for _, m := range matches {
		tag := strings.ToLower(m[1])
		unique[tag] = struct{}{}
	}
	result := make([]string, 0, len(unique))
	for tag := range unique {
		result = append(result, tag)
	}
	return result
}

func (s *PostService) CreatePost(userID string, req *model.CreatePostRequest) (*model.Post, error) {
	post := &model.Post{
		ID: uuid.New().String(),
		UserID: userID,
		Content: req.Content,
		ImageURLs: model.ImageURLs(req.ImageURLs),
	}

	err := s.postRepo.Create(post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Xử lý hashtag
	hashtags := extractHashtags(req.Content)
	if len(hashtags) > 0 {
		err = s.postRepo.AddHashtagsToPost(post.ID, hashtags)
		if err != nil {
			return nil, fmt.Errorf("failed to add hashtags: %w", err)
		}
	}

	return s.postRepo.GetByID(post.ID, &userID)
}

func (s *PostService) GetPost(postID string, userID *string) (*model.Post, error) {
	post, err := s.postRepo.GetByID(postID, userID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

func (s *PostService) UpdatePost(postID, userID string, req *model.UpdatePostRequest) (*model.Post, error) {
	// First check if the post exists and belongs to the user
	existingPost, err := s.postRepo.GetByID(postID, &userID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	if existingPost.UserID != userID {
		return nil, fmt.Errorf("you can only edit your own posts")
	}

	post := &model.Post{
		ID:        postID,
		UserID:    userID,
		Content:   req.Content,
		ImageURLs: model.ImageURLs(req.ImageURLs),
	}

	err = s.postRepo.Update(post)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Xử lý hashtag (cập nhật lại toàn bộ hashtag cho post)
	hashtags := extractHashtags(req.Content)
	if len(hashtags) > 0 {
		err = s.postRepo.ClearPostHashtags(postID)
		if err != nil {
			return nil, fmt.Errorf("failed to clear old hashtags: %w", err)
		}
		err = s.postRepo.AddHashtagsToPost(postID, hashtags)
		if err != nil {
			return nil, fmt.Errorf("failed to add hashtags: %w", err)
		}
	} else {
		err = s.postRepo.ClearPostHashtags(postID)
		if err != nil {
			return nil, fmt.Errorf("failed to clear old hashtags: %w", err)
		}
	}

	return s.postRepo.GetByID(postID, &userID)
}

func (s *PostService) DeletePost(postID, userID string) error {
	err := s.postRepo.Delete(postID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (s *PostService) GetFeed(userID string, pagination *utils.PaginationParams) (*model.PostsResponse, error) {
	paginationResult := pagination.Calculate()

	posts, totalCount, err := s.postRepo.GetFeed(userID, paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.PostsResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       paginationResult.Page,
		PageSize:   paginationResult.PageSize,
		HasMore:    hasMore,
	}, nil
}

func (s *PostService) GetUserPosts(userID string, viewerID *string, pagination *utils.PaginationParams) (*model.PostsResponse, error) {
	paginationResult := pagination.Calculate()

	posts, totalCount, err := s.postRepo.GetUserPosts(userID, viewerID, paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.PostsResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       paginationResult.Page,
		PageSize:   paginationResult.PageSize,
		HasMore:    hasMore,
	}, nil
}

func (s *PostService) ToggleLike(postID, userID string) (bool, error) {
	// Check if post is already liked
	isLiked, err := s.postRepo.IsPostLikedByUser(postID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}

	if isLiked {
		err = s.postRepo.UnlikePost(postID, userID)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		err = s.postRepo.LikePost(postID, userID)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func (s *PostService) LikePost(postID, userID string) error {
	return s.postRepo.LikePost(postID, userID)
}

func (s *PostService) UnlikePost(postID, userID string) error {
	return s.postRepo.UnlikePost(postID, userID)
}

func (s *PostService) IsPostLikedByUser(postID, userID string) (bool, error) {
	return s.postRepo.IsPostLikedByUser(postID, userID)
}

func (s *PostService) GetExplorePosts(userID *string, pagination *utils.PaginationParams) (*model.PostsResponse, error) {
	paginationResult := pagination.Calculate()
	posts, totalCount, err := s.postRepo.GetUserPosts("", userID, paginationResult) // userID rỗng để lấy tất cả
	if err != nil {
		return nil, fmt.Errorf("failed to get explore posts: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.PostsResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       paginationResult.Page,
		PageSize:   paginationResult.PageSize,
		HasMore:    hasMore,
	}, nil
}

func (s *PostService) GetPostStats(postID string) (map[string]interface{}, error) {
	post, err := s.postRepo.GetByID(postID, nil)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	stats := map[string]interface{}{
		"like_count":    post.LikeCount,
		"comment_count": post.CommentCount,
		"created_at":    post.CreatedAt,
		"has_images":    len(post.ImageURLs) > 0,
		"image_count":   len(post.ImageURLs),
	}

	return stats, nil
}

func (s *PostService) SearchPosts(query string, page, pageSize int) (*model.PostsResponse, error) {
	posts, totalCount, err := s.postRepo.SearchPosts(query, page, pageSize)
	if err != nil {
		return nil, err
	}
	hasMore := utils.CalculateHasMore(totalCount, page, pageSize)
	return &model.PostsResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    hasMore,
	}, nil
}

// SearchPostsByContent chỉ theo content
func (s *PostService) SearchPostsByContent(query string, page, pageSize int) (*model.PostsResponse, error) {
	posts, totalCount, err := s.postRepo.SearchPostsByContent(query, page, pageSize)
	if err != nil {
		return nil, err
	}
	hasMore := utils.CalculateHasMore(totalCount, page, pageSize)
	return &model.PostsResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    hasMore,
	}, nil
}

// Lấy danh sách hashtag của post
func (s *PostService) GetHashtagsByPost(postID string) ([]model.Hashtag, error) {
	return s.postRepo.GetHashtagsByPost(postID)
}

// Lấy danh sách post theo hashtag
func (s *PostService) GetPostsByHashtag(hashtag string, limit, offset int) ([]model.Post, error) {
	return s.postRepo.GetPostsByHashtag(hashtag, limit, offset)
}

func (s *PostService) SearchHashtags(query string, page, pageSize int) ([]model.Hashtag, int64, bool, error) {
	hashtags, totalCount, err := s.postRepo.SearchHashtags(query, page, pageSize)
	if err != nil {
		return nil, 0, false, err
	}
	hasMore := utils.CalculateHasMore(totalCount, page, pageSize)
	return hashtags, totalCount, hasMore, nil
}
