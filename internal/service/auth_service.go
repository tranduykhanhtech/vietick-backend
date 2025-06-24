package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"vietick-backend/internal/model"
	"vietick-backend/internal/repository"
	"vietick-backend/internal/utils"
	"vietick-backend/pkg/email"
	"vietick-backend/pkg/jwt"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	authRepo    *repository.AuthRepository
	jwtManager  *jwt.JWTManager
	emailService *email.EmailService
}

func NewAuthService(userRepo *repository.UserRepository, authRepo *repository.AuthRepository, 
	jwtManager *jwt.JWTManager, emailService *email.EmailService) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		jwtManager:  jwtManager,
		emailService: emailService,
	}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	existingUser, _ = s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this username already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate email verification token
	verificationToken, err := utils.GenerateEmailVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create user
	user := &model.User{
		ID:                        uuid.New().String(),
		Username:                  req.Username,
		Email:                     req.Email,
		PasswordHash:              passwordHash,
		FullName:                  req.FullName,
		IsEmailVerified:           false,
		EmailVerificationToken:    &verificationToken,
		EmailVerificationExpiresAt: timePtr(time.Now().Add(24 * time.Hour)),
		IdentityVerificationStatus: model.IdentityVerificationNone,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send verification email
	err = s.emailService.SendEmailVerification(user.Email, user.FullName, verificationToken)
	if err != nil {
		// Log error but don't fail registration
		fmt.Printf("Failed to send verification email: %v\n", err)
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token
	err = s.storeRefreshToken(user.ID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Get user profile
	profile, err := s.userRepo.GetProfile(user.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *profile,
		ExpiresIn:    int64(s.jwtManager.GetAccessTokenExpiry().Seconds()),
	}, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token
	err = s.storeRefreshToken(user.ID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Get user profile
	profile, err := s.userRepo.GetProfile(user.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *profile,
		ExpiresIn:    int64(s.jwtManager.GetAccessTokenExpiry().Seconds()),
	}, nil
}

func (s *AuthService) RefreshToken(req *model.RefreshTokenRequest) (*model.AuthResponse, error) {
	// Validate refresh token
	_, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token exists in database
	tokenHash := s.hashToken(req.RefreshToken)
	storedToken, err := s.authRepo.GetRefreshToken(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Remove old refresh token and store new one
	err = s.authRepo.DeleteRefreshToken(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	err = s.storeRefreshToken(user.ID, newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	// Get user profile
	profile, err := s.userRepo.GetProfile(user.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         *profile,
		ExpiresIn:    int64(s.jwtManager.GetAccessTokenExpiry().Seconds()),
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := s.hashToken(refreshToken)
	return s.authRepo.DeleteRefreshToken(tokenHash)
}

func (s *AuthService) LogoutAll(userID string) error {
	return s.authRepo.DeleteUserRefreshTokens(userID)
}

func (s *AuthService) VerifyEmail(req *model.VerifyEmailRequest) error {
	user, err := s.userRepo.GetByEmailVerificationToken(req.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	if user.IsEmailVerified {
		return fmt.Errorf("email already verified")
	}

	err = s.userRepo.VerifyEmail(user.ID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

func (s *AuthService) ResendEmailVerification(userID string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.IsEmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Generate new verification token
	verificationToken, err := utils.GenerateEmailVerificationToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	user.EmailVerificationToken = &verificationToken
	expiresAt := time.Now().Add(24 * time.Hour)
	user.EmailVerificationExpiresAt = &expiresAt

	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Send verification email
	err = s.emailService.SendEmailVerification(user.Email, user.FullName, verificationToken)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (s *AuthService) ChangePassword(userID string, req *model.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check current password
	if !utils.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	err = s.userRepo.UpdatePassword(user.ID, newPasswordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Logout from all devices for security
	err = s.authRepo.DeleteUserRefreshTokens(user.ID)
	if err != nil {
		return fmt.Errorf("failed to logout from all devices: %w", err)
	}

	return nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*jwt.Claims, error) {
	return s.jwtManager.ValidateAccessToken(tokenString)
}

func (s *AuthService) storeRefreshToken(userID string, refreshToken string) error {
	// Check if user has too many refresh tokens (limit to 5 devices)
	count, err := s.authRepo.GetUserRefreshTokenCount(userID)
	if err != nil {
		return err
	}

	if count >= 5 {
		// Remove oldest token
		err = s.authRepo.DeleteOldestUserRefreshToken(userID)
		if err != nil {
			return err
		}
	}

	refreshTokenModel := &model.RefreshToken{
		ID: uuid.New().String(),
		UserID: userID,
		TokenHash: s.hashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenExpiry()),
	}

	return s.authRepo.CreateRefreshToken(refreshTokenModel)
}

func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) CleanupExpiredTokens() error {
	return s.authRepo.CleanupExpiredTokens()
}

func timePtr(t time.Time) *time.Time {
	return &t
}
