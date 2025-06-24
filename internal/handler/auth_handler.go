package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vietick-backend/internal/middleware"
	"vietick-backend/internal/model"
	"vietick-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Registration data"
// @Success 201 {object} model.AuthResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, authResponse)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	authResponse, err := h.authService.RefreshToken(&req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// LogoutAll godoc
// @Summary Logout from all devices
// @Description Invalidate all refresh tokens for the user
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	err := h.authService.LogoutAll(userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out from all devices successfully",
	})
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user's email address using verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.VerifyEmailRequest true "Email verification token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req model.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	err := h.authService.VerifyEmail(&req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// ResendEmailVerification godoc
// @Summary Resend email verification
// @Description Send a new email verification link to the user
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /auth/resend-verification [post]
func (h *AuthHandler) ResendEmailVerification(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	err := h.authService.ResendEmailVerification(userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully. Please login again.",
	})
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} model.UserProfile
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// This would typically be handled by UserHandler, but including here for convenience
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "User authenticated",
	})
}

// CheckToken godoc
// @Summary Check token validity
// @Description Check if the provided access token is valid
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /auth/check [get]
func (h *AuthHandler) CheckToken(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	claims, _ := middleware.GetClaims(c)

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"user_id":  userID,
		"username": claims.Username,
		"email":    claims.Email,
	})
}
