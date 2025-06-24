package service

import (
	"fmt"

	"vietick-backend/internal/model"
	"vietick-backend/internal/repository"
	"vietick-backend/internal/utils"
	"github.com/google/uuid"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
}

func NewCommentService(commentRepo *repository.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (s *CommentService) CreateComment(userID, postID string, req *model.CreateCommentRequest) (*model.Comment, error) {
	comment := &model.Comment{
		ID: uuid.New().String(),
		PostID: postID,
		UserID: userID,
		Content: req.Content,
	}

	err := s.commentRepo.Create(comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Get the complete comment with user information
	return s.commentRepo.GetByID(comment.ID, &userID)
}

func (s *CommentService) GetComment(commentID string, userID *string) (*model.Comment, error) {
	comment, err := s.commentRepo.GetByID(commentID, userID)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	return comment, nil
}

func (s *CommentService) GetPostComments(postID string, userID *string, pagination *utils.PaginationParams) (*model.CommentsResponse, error) {
	paginationResult := pagination.Calculate()

	comments, totalCount, err := s.commentRepo.GetPostComments(postID, userID, paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.CommentsResponse{
		Comments:   comments,
		TotalCount: totalCount,
		Page:       paginationResult.Page,
		PageSize:   paginationResult.PageSize,
		HasMore:    hasMore,
	}, nil
}

func (s *CommentService) UpdateComment(commentID, userID string, req *model.CreateCommentRequest) (*model.Comment, error) {
	// First check if the comment exists and belongs to the user
	existingComment, err := s.commentRepo.GetByID(commentID, &userID)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	if existingComment.UserID != userID {
		return nil, fmt.Errorf("you can only edit your own comments")
	}

	comment := &model.Comment{
		ID:      commentID,
		UserID:  userID,
		Content: req.Content,
	}

	err = s.commentRepo.Update(comment)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	// Get the updated comment
	return s.commentRepo.GetByID(commentID, &userID)
}

func (s *CommentService) DeleteComment(commentID, userID string) error {
	err := s.commentRepo.Delete(commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (s *CommentService) LikeComment(commentID, userID string) error {
	err := s.commentRepo.LikeComment(commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to like comment: %w", err)
	}

	return nil
}

func (s *CommentService) UnlikeComment(commentID, userID string) error {
	err := s.commentRepo.UnlikeComment(commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to unlike comment: %w", err)
	}

	return nil
}

func (s *CommentService) ToggleLike(commentID, userID string) (bool, error) {
	// Check if comment is already liked
	isLiked, err := s.commentRepo.IsCommentLikedByUser(commentID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}

	if isLiked {
		err = s.UnlikeComment(commentID, userID)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		err = s.LikeComment(commentID, userID)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func (s *CommentService) IsCommentLikedByUser(commentID, userID string) (bool, error) {
	return s.commentRepo.IsCommentLikedByUser(commentID, userID)
}

func (s *CommentService) GetCommentStats(commentID string) (map[string]interface{}, error) {
	comment, err := s.commentRepo.GetByID(commentID, nil)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	stats := map[string]interface{}{
		"like_count": comment.LikeCount,
		"created_at": comment.CreatedAt,
		"updated_at": comment.UpdatedAt,
	}

	return stats, nil
}
