package service

import (
	"fmt"

	"vietick-backend/internal/model"
	"vietick-backend/internal/repository"
	"vietick-backend/internal/utils"
)

type FollowService struct {
	followRepo *repository.FollowRepository
}

func NewFollowService(followRepo *repository.FollowRepository) *FollowService {
	return &FollowService{
		followRepo: followRepo,
	}
}

func (s *FollowService) Follow(followerID, followingID string) error {
	err := s.followRepo.Follow(followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}
	return nil
}

func (s *FollowService) Unfollow(followerID, followingID string) error {
	err := s.followRepo.Unfollow(followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	return nil
}

func (s *FollowService) ToggleFollow(followerID, followingID string) (bool, error) {
	// Check if already following
	isFollowing, err := s.followRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}

	if isFollowing {
		err = s.Unfollow(followerID, followingID)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		err = s.Follow(followerID, followingID)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

func (s *FollowService) IsFollowing(followerID, followingID string) (bool, error) {
	return s.followRepo.IsFollowing(followerID, followingID)
}

func (s *FollowService) GetFollowers(userID string, viewerID *string, pagination *utils.PaginationParams) (*model.FollowersResponse, error) {
	paginationResult := pagination.Calculate()

	users, totalCount, err := s.followRepo.GetFollowers(userID, viewerID, paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.FollowersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       paginationResult.Page,
		PageSize:   paginationResult.PageSize,
		HasMore:    hasMore,
	}, nil
}

func (s *FollowService) GetFollowing(userID string, viewerID *string, pagination *utils.PaginationParams) (*model.FollowingResponse, error) {
	paginationResult := pagination.Calculate()

	users, totalCount, err := s.followRepo.GetFollowing(userID, viewerID, paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.FollowingResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       paginationResult.Page,
		PageSize:   paginationResult.PageSize,
		HasMore:    hasMore,
	}, nil
}

func (s *FollowService) GetFollowCounts(userID string) (followersCount, followingCount int64, err error) {
	return s.followRepo.GetFollowCounts(userID)
}

func (s *FollowService) GetMutualFollows(userID1, userID2 string, limit int) ([]model.UserProfile, error) {
	users, err := s.followRepo.GetMutualFollows(userID1, userID2, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get mutual follows: %w", err)
	}

	return users, nil
}

func (s *FollowService) GetFollowRelationship(userID1, userID2 string) (map[string]bool, error) {
	isFollowing, err := s.followRepo.IsFollowing(userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user1 follows user2: %w", err)
	}

	isFollowedBy, err := s.followRepo.IsFollowing(userID2, userID1)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user2 follows user1: %w", err)
	}

	relationship := map[string]bool{
		"is_following":   isFollowing,
		"is_followed_by": isFollowedBy,
		"is_mutual":      isFollowing && isFollowedBy,
	}

	return relationship, nil
}

func (s *FollowService) GetFollowStats(userID string) (map[string]interface{}, error) {
	followersCount, followingCount, err := s.followRepo.GetFollowCounts(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get follow counts: %w", err)
	}

	stats := map[string]interface{}{
		"followers_count": followersCount,
		"following_count": followingCount,
		"follow_ratio":    calculateFollowRatio(followersCount, followingCount),
	}

	return stats, nil
}

func calculateFollowRatio(followers, following int64) float64 {
	if following == 0 {
		if followers == 0 {
			return 0.0
		}
		return float64(followers) // Infinite ratio, return followers count
	}
	return float64(followers) / float64(following)
}

func (s *FollowService) GetRecommendedUsers(userID string, limit int) ([]model.UserProfile, error) {
	// Get users followed by people the current user follows
	// This is a "friends of friends" recommendation
	
	// For now, return empty slice as this would require a complex query
	// In a real implementation, you might:
	// 1. Get users followed by users that the current user follows
	// 2. Filter out users already followed
	// 3. Sort by mutual connections or other relevance metrics
	
	return []model.UserProfile{}, nil
}

func (s *FollowService) BulkFollow(followerID string, followingIDs []string) ([]string, []string, error) {
	var successful []string
	var errors []string

	for _, followingID := range followingIDs {
		err := s.Follow(followerID, followingID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to follow user %s: %v", followingID, err))
		} else {
			successful = append(successful, followingID)
		}
	}

	return successful, errors, nil
}

func (s *FollowService) BulkUnfollow(followerID string, followingIDs []string) ([]string, []string, error) {
	var successful []string
	var errors []string

	for _, followingID := range followingIDs {
		err := s.Unfollow(followerID, followingID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to unfollow user %s: %v", followingID, err))
		} else {
			successful = append(successful, followingID)
		}
	}

	return successful, errors, nil
}
