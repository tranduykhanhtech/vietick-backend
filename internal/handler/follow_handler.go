package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"vietick-backend/internal/middleware"
	"vietick-backend/internal/service"
	"vietick-backend/internal/utils"
)

type FollowHandler struct {
	followService *service.FollowService
}

func NewFollowHandler(followService *service.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

// Follow godoc
// @Summary Follow a user
// @Description Follow a user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID to follow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/follow [post]
func (h *FollowHandler) Follow(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	followingID := c.Param("user_id")
	if followingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err := h.followService.Follow(followerID, followingID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "User followed successfully",
		"following":    true,
		"follower_id":  followerID,
		"following_id": followingID,
	})
}

// Unfollow godoc
// @Summary Unfollow a user
// @Description Unfollow a user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID to unfollow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/unfollow [post]
func (h *FollowHandler) Unfollow(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	followingID := c.Param("user_id")
	if followingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err := h.followService.Unfollow(followerID, followingID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "User unfollowed successfully",
		"following":    false,
		"follower_id":  followerID,
		"following_id": followingID,
	})
}

// ToggleFollow godoc
// @Summary Toggle follow status
// @Description Follow or unfollow a user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID to toggle follow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/toggle-follow [post]
func (h *FollowHandler) ToggleFollow(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	followingID := c.Param("user_id")
	if followingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	following, err := h.followService.ToggleFollow(followerID, followingID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	message := "User unfollowed successfully"
	if following {
		message = "User followed successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"following":    following,
		"follower_id":  followerID,
		"following_id": followingID,
	})
}

// GetFollowers godoc
// @Summary Get user followers
// @Description Get a list of users who follow the specified user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.FollowersResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/followers [get]
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	viewerID := middleware.GetUserIDPtr(c)

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.followService.GetFollowers(userID, viewerID, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFollowing godoc
// @Summary Get users followed by user
// @Description Get a list of users that the specified user follows
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.FollowingResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/following [get]
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	viewerID := middleware.GetUserIDPtr(c)

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.followService.GetFollowing(userID, viewerID, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFollowStatus godoc
// @Summary Get follow status
// @Description Check if the authenticated user follows the specified user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID to check"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/follow-status [get]
func (h *FollowHandler) GetFollowStatus(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	followingID := c.Param("user_id")
	if followingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	isFollowing, err := h.followService.IsFollowing(followerID, followingID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_following": isFollowing,
		"follower_id":  followerID,
		"following_id": followingID,
	})
}

// GetFollowCounts godoc
// @Summary Get follow counts
// @Description Get follower and following counts for a user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/follow-counts [get]
func (h *FollowHandler) GetFollowCounts(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followersCount, followingCount, err := h.followService.GetFollowCounts(userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":         userID,
		"followers_count": followersCount,
		"following_count": followingCount,
	})
}

// GetMutualFollows godoc
// @Summary Get mutual follows
// @Description Get users that both the authenticated user and specified user follow
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} []model.UserProfile
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/mutual-follows [get]
func (h *FollowHandler) GetMutualFollows(c *gin.Context) {
	currentUserID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	users, err := h.followService.GetMutualFollows(currentUserID, userID, limit)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GetFollowRelationship godoc
// @Summary Get follow relationship
// @Description Get the follow relationship between authenticated user and specified user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/relationship [get]
func (h *FollowHandler) GetFollowRelationship(c *gin.Context) {
	currentUserID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	relationship, err := h.followService.GetFollowRelationship(currentUserID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, relationship)
}

// GetFollowStats godoc
// @Summary Get follow statistics
// @Description Get follow statistics for a user
// @Tags follows
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/follow-stats [get]
func (h *FollowHandler) GetFollowStats(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.followService.GetFollowStats(userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// BulkFollow godoc
// @Summary Bulk follow users
// @Description Follow multiple users at once
// @Tags follows
// @Accept json
// @Produce json
// @Param request body map[string][]string true "User IDs to follow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /follows/bulk-follow [post]
func (h *FollowHandler) BulkFollow(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		UserIDs []string `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	successful, errors, err := h.followService.BulkFollow(followerID, req.UserIDs)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"successful_follows": successful,
		"successful_count":   len(successful),
		"errors":            errors,
		"error_count":       len(errors),
	})
}

// BulkUnfollow godoc
// @Summary Bulk unfollow users
// @Description Unfollow multiple users at once
// @Tags follows
// @Accept json
// @Produce json
// @Param request body map[string][]string true "User IDs to unfollow"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /follows/bulk-unfollow [post]
func (h *FollowHandler) BulkUnfollow(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		UserIDs []string `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	successful, errors, err := h.followService.BulkUnfollow(followerID, req.UserIDs)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"successful_unfollows": successful,
		"successful_count":     len(successful),
		"errors":              errors,
		"error_count":         len(errors),
	})
}
