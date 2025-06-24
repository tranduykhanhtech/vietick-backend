package service

import (
	"fmt"

	"vietick-backend/internal/model"
	"vietick-backend/internal/repository"
	"vietick-backend/internal/utils"
	"vietick-backend/pkg/email"
	"github.com/google/uuid"
)

type VerificationService struct {
	verificationRepo *repository.VerificationRepository
	userRepo         *repository.UserRepository
	emailService     *email.EmailService
}

func NewVerificationService(verificationRepo *repository.VerificationRepository, 
	userRepo *repository.UserRepository, emailService *email.EmailService) *VerificationService {
	return &VerificationService{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
		emailService:     emailService,
	}
}

func (s *VerificationService) SubmitIdentityVerification(userID string, req *model.SubmitIdentityVerificationRequest) (*model.IdentityVerification, error) {
	// Check if user already has a pending verification
	hasPending, err := s.verificationRepo.HasPendingVerification(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check pending verification: %w", err)
	}

	if hasPending {
		return nil, fmt.Errorf("you already have a pending verification request")
	}

	// Check if user is already verified
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.IsVerified {
		return nil, fmt.Errorf("user is already verified")
	}

	// Create verification request
	verification := &model.IdentityVerification{
		ID: uuid.New().String(),
		UserID: userID,
		FullName: req.FullName,
		IDNumber: req.IDNumber,
		IDType: req.IDType,
		FrontImageURL: req.FrontImageURL,
		BackImageURL: req.BackImageURL,
		SelfieImageURL: req.SelfieImageURL,
		Status: model.IdentityVerificationPending,
	}

	err = s.verificationRepo.Create(verification)
	if err != nil {
		return nil, fmt.Errorf("failed to submit verification: %w", err)
	}

	// Update user verification status
	user.IdentityVerificationStatus = model.IdentityVerificationPending
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user verification status: %w", err)
	}

	// Get the complete verification with user information
	return s.verificationRepo.GetByID(verification.ID)
}

func (s *VerificationService) GetUserVerification(userID string) (*model.IdentityVerification, error) {
	verification, err := s.verificationRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("verification not found")
	}

	return verification, nil
}

func (s *VerificationService) GetVerification(verificationID string) (*model.IdentityVerification, error) {
	verification, err := s.verificationRepo.GetByID(verificationID)
	if err != nil {
		return nil, fmt.Errorf("verification not found")
	}

	return verification, nil
}

func (s *VerificationService) GetPendingVerifications(pagination *utils.PaginationParams) (*model.IdentityVerificationsResponse, error) {
	paginationResult := pagination.Calculate()

	verifications, totalCount, err := s.verificationRepo.GetPending(paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending verifications: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.IdentityVerificationsResponse{
		Verifications: verifications,
		TotalCount:    totalCount,
		Page:          paginationResult.Page,
		PageSize:      paginationResult.PageSize,
		HasMore:       hasMore,
	}, nil
}

func (s *VerificationService) GetAllVerifications(status *model.IdentityVerificationStatus, pagination *utils.PaginationParams) (*model.IdentityVerificationsResponse, error) {
	paginationResult := pagination.Calculate()

	verifications, totalCount, err := s.verificationRepo.GetAll(status, paginationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get verifications: %w", err)
	}

	hasMore := utils.CalculateHasMore(totalCount, paginationResult.Page, paginationResult.PageSize)

	return &model.IdentityVerificationsResponse{
		Verifications: verifications,
		TotalCount:    totalCount,
		Page:          paginationResult.Page,
		PageSize:      paginationResult.PageSize,
		HasMore:       hasMore,
	}, nil
}

func (s *VerificationService) ReviewVerification(verificationID string, reviewedBy string, req *model.ReviewIdentityVerificationRequest) (*model.IdentityVerification, error) {
	// Get verification
	verification, err := s.verificationRepo.GetByID(verificationID)
	if err != nil {
		return nil, fmt.Errorf("verification not found")
	}

	if verification.Status != model.IdentityVerificationPending {
		return nil, fmt.Errorf("verification has already been reviewed")
	}

	// Review verification
	err = s.verificationRepo.Review(verificationID, req.Status, req.AdminNotes, reviewedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to review verification: %w", err)
	}

	// If approved, send approval email
	if req.Status == model.IdentityVerificationApproved {
		user, err := s.userRepo.GetByID(verification.UserID)
		if err == nil { // Don't fail the review if email fails
			err = s.emailService.SendVerificationApproval(user.Email, user.FullName)
			if err != nil {
				// Log error but don't fail the operation
				fmt.Printf("Failed to send verification approval email: %v\n", err)
			}
		}
	}

	// Get updated verification
	return s.verificationRepo.GetByID(verificationID)
}

func (s *VerificationService) GetVerificationStats() (map[string]interface{}, error) {
	// Get counts for different statuses
	pendingStatus := model.IdentityVerificationPending
	approvedStatus := model.IdentityVerificationApproved
	rejectedStatus := model.IdentityVerificationRejected

	pendingPagination := utils.PaginationResult{Offset: 0, Limit: 1, Page: 1, PageSize: 1}

	pendingVerifications, pendingCount, err := s.verificationRepo.GetAll(&pendingStatus, pendingPagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending count: %w", err)
	}
	_ = pendingVerifications // We only need the count

	approvedVerifications, approvedCount, err := s.verificationRepo.GetAll(&approvedStatus, pendingPagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved count: %w", err)
	}
	_ = approvedVerifications // We only need the count

	rejectedVerifications, rejectedCount, err := s.verificationRepo.GetAll(&rejectedStatus, pendingPagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejected count: %w", err)
	}
	_ = rejectedVerifications // We only need the count

	totalCount := pendingCount + approvedCount + rejectedCount

	stats := map[string]interface{}{
		"total_submissions":    totalCount,
		"pending_submissions":  pendingCount,
		"approved_submissions": approvedCount,
		"rejected_submissions": rejectedCount,
		"approval_rate":        calculateApprovalRate(approvedCount, totalCount),
	}

	return stats, nil
}

func calculateApprovalRate(approved, total int64) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(approved) / float64(total) * 100
}

func (s *VerificationService) CanSubmitVerification(userID string) (bool, string, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, "User not found", err
	}

	// Check if already verified
	if user.IsVerified {
		return false, "User is already verified", nil
	}

	// Check if email is verified
	if !user.IsEmailVerified {
		return false, "Email must be verified before submitting identity verification", nil
	}

	// Check if has pending verification
	hasPending, err := s.verificationRepo.HasPendingVerification(userID)
	if err != nil {
		return false, "Failed to check pending verification", err
	}

	if hasPending {
		return false, "A verification request is already pending", nil
	}

	return true, "Can submit verification", nil
}

func (s *VerificationService) DeleteVerification(verificationID string) error {
	err := s.verificationRepo.Delete(verificationID)
	if err != nil {
		return fmt.Errorf("failed to delete verification: %w", err)
	}

	return nil
}

func (s *VerificationService) GetVerifiedUsers(pagination *utils.PaginationParams) (*model.FollowersResponse, error) {
	// This would need a custom repository method to get verified users
	// For now, return empty response as placeholder
	return &model.FollowersResponse{
		Users:      []model.UserProfile{},
		TotalCount: 0,
		Page:       pagination.Calculate().Page,
		PageSize:   pagination.Calculate().PageSize,
		HasMore:    false,
	}, nil
}
