package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"vietick-backend/internal/middleware"
	"vietick-backend/internal/model"
	"vietick-backend/internal/service"
	"vietick-backend/internal/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get user profile by user ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.UserProfile
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	viewerID := middleware.GetUserIDPtr(c)

	profile, err := h.userService.GetProfile(userID, viewerID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetProfileByUsername godoc
// @Summary Get user profile by username
// @Description Get user profile by username
// @Tags users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} model.UserProfile
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/username/{username} [get]
func (h *UserHandler) GetProfileByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	viewerID := middleware.GetUserIDPtr(c)

	profile, err := h.userService.GetProfileByUsername(username, viewerID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetCurrentProfile godoc
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} model.UserProfile
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) GetCurrentProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profile, err := h.userService.GetProfile(userID, &userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} model.UserProfile
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	profile, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users by username or full name
// @Tags users
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.FollowersResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	viewerID := middleware.GetUserIDPtr(c)

	response, err := h.userService.SearchUsers(query, &pagination, viewerID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get statistics for a user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.userService.GetUserStats(userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// CheckUsernameAvailability godoc
// @Summary Check username availability
// @Description Check if a username is available
// @Tags users
// @Produce json
// @Param username query string true "Username to check"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} middleware.ErrorResponse
// @Router /users/check-username [get]
func (h *UserHandler) CheckUsernameAvailability(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username is required",
		})
		return
	}

	available, err := h.userService.CheckUsernameAvailability(username)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"username":  username,
	})
}

// CheckEmailAvailability godoc
// @Summary Check email availability
// @Description Check if an email is available
// @Tags users
// @Produce json
// @Param email query string true "Email to check"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} middleware.ErrorResponse
// @Router /users/check-email [get]
func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})
		return
	}

	available, err := h.userService.CheckEmailAvailability(email)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"email":     email,
	})
}

// GetRecommendedUsers godoc
// @Summary Get recommended users
// @Description Get users recommended for the authenticated user
// @Tags users
// @Produce json
// @Param limit query int false "Number of recommendations" default(10)
// @Success 200 {object} []model.UserProfile
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/recommended [get]
func (h *UserHandler) GetRecommendedUsers(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	users, err := h.userService.GetRecommendedUsers(userID, limit)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// UpdateUsername godoc
// @Summary Update username
// @Description Update the authenticated user's username
// @Tags users
// @Accept json
// @Produce json
// @Param request body map[string]string true "New username"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/me/username [put]
func (h *UserHandler) UpdateUsername(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	err := h.userService.UpdateUsername(userID, req.Username)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Username updated successfully",
		"username": req.Username,
	})
}

// UpdateEmail godoc
// @Summary Update email
// @Description Update the authenticated user's email
// @Tags users
// @Accept json
// @Produce json
// @Param request body map[string]string true "New email"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/me/email [put]
func (h *UserHandler) UpdateEmail(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	err := h.userService.UpdateEmail(userID, req.Email)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email updated successfully. Please verify your new email address.",
		"email":   req.Email,
	})
}
