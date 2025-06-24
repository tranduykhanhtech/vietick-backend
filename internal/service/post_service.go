package service

import (
	"fmt"

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

	// Get the complete post with user information
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

	// Get the updated post
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

func (s *PostService) SearchPosts(query string, userID *string, pagination *utils.PaginationParams) (*model.PostsResponse, error) {
	return &model.PostsResponse{
		Posts:      []model.Post{},
		TotalCount: 0,
		Page:       pagination.Calculate().Page,
		PageSize:   pagination.Calculate().PageSize,
		HasMore:    false,
	}, nil
}
