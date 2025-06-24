package service

import (
	"fmt"
	"time"

	"vietick-backend/internal/model"
	"vietick-backend/internal/repository"
	"vietick-backend/internal/utils"
)

type UserService struct {
	userRepo   *repository.UserRepository
	followRepo *repository.FollowRepository
}

func NewUserService(userRepo *repository.UserRepository, followRepo *repository.FollowRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *UserService) GetProfile(userID string, viewerID *string) (*model.UserProfile, error) {
	profile, err := s.userRepo.GetProfile(userID, viewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return profile, nil
}

func (s *UserService) GetProfileByUsername(username string, viewerID *string) (*model.UserProfile, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.GetProfile(user.ID, viewerID)
}

func (s *UserService) UpdateProfile(userID string, req *model.UpdateProfileRequest) (*model.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Update fields if provided
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return s.GetProfile(userID, nil)
}

func (s *UserService) SearchUsers(query string, pagination *utils.PaginationParams, viewerID *string) (*model.FollowersResponse, error) {
	paginationResult := pagination.Calculate()
	
	users, totalCount, err := s.userRepo.SearchUsers(query, paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	// If viewer is provided, get follow relationships
	if viewerID != nil {
		for i := range users {
			isFollowing, _ := s.followRepo.IsFollowing(*viewerID, users[i].ID)
			users[i].IsFollowing = isFollowing
			
			isFollowedBy, _ := s.followRepo.IsFollowing(users[i].ID, *viewerID)
			users[i].IsFollowedBy = isFollowedBy
		}
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

func (s *UserService) GetUserByID(userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) CheckUsernameAvailability(username string) (bool, error) {
	_, err := s.userRepo.GetByUsername(username)
	if err != nil {
		// If user not found, username is available
		return true, nil
	}

	// If user found, username is not available
	return false, nil
}

func (s *UserService) CheckEmailAvailability(email string) (bool, error) {
	_, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// If user not found, email is available
		return true, nil
	}

	// If user found, email is not available
	return false, nil
}

func (s *UserService) UpdateUsername(userID string, newUsername string) error {
	// Check if username is available
	available, err := s.CheckUsernameAvailability(newUsername)
	if err != nil {
		return fmt.Errorf("failed to check username availability: %w", err)
	}

	if !available {
		return fmt.Errorf("username is already taken")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.Username = newUsername
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

func (s *UserService) UpdateEmail(userID string, newEmail string) error {
	// Check if email is available
	available, err := s.CheckEmailAvailability(newEmail)
	if err != nil {
		return fmt.Errorf("failed to check email availability: %w", err)
	}

	if !available {
		return fmt.Errorf("email is already in use")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Generate new email verification token
	verificationToken, err := utils.GenerateEmailVerificationToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	user.Email = newEmail
	user.IsEmailVerified = false
	user.EmailVerificationToken = &verificationToken
	expiresAt := time.Now().Add(24 * time.Hour)
	user.EmailVerificationExpiresAt = &expiresAt

	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

func (s *UserService) GetRecommendedUsers(userID string, limit int) ([]model.UserProfile, error) {
	// This is a simple implementation. In a real application, you might want to
	// implement more sophisticated recommendation algorithms based on:
	// - Mutual follows
	// - Similar interests
	// - Popular users
	// - Geographic location
	// - Activity patterns

	pagination := utils.PaginationResult{
		Offset:   0,
		Limit:    limit,
		Page:     1,
		PageSize: limit,
	}

	// For now, get recent verified users that the current user doesn't follow
	users, _, err := s.userRepo.SearchUsers("", pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended users: %w", err)
	}

	// Filter out current user and users already followed
	var recommendations []model.UserProfile

	for _, user := range users {
		if user.ID == userID {
			continue
		}

		isFollowing, _ := s.followRepo.IsFollowing(userID, user.ID)
		if isFollowing {
			continue
		}

		user.IsFollowing = false
		recommendations = append(recommendations, user)

		if len(recommendations) >= limit {
			break
		}
	}

	return recommendations, nil
}

func (s *UserService) GetUserStats(userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	profile, err := s.userRepo.GetProfile(userID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	stats := map[string]interface{}{
		"posts_count":     profile.PostsCount,
		"followers_count": profile.FollowersCount,
		"following_count": profile.FollowingCount,
		"is_verified":     profile.IsVerified,
		"is_email_verified": user.IsEmailVerified,
		"account_age_days": int(time.Since(user.CreatedAt).Hours() / 24),
		"verification_status": user.IdentityVerificationStatus,
	}

	return stats, nil
}
