package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vietick-backend/internal/middleware"
	"vietick-backend/internal/model"
	"vietick-backend/internal/service"
	"vietick-backend/internal/utils"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body model.CreatePostRequest true "Post data"
// @Success 201 {object} model.Post
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	post, err := h.postService.CreatePost(userID, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPost godoc
// @Summary Get a post
// @Description Get a post by ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} model.Post
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userID := middleware.GetUserIDPtr(c)

	post, err := h.postService.GetPost(postID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update a post (only by owner)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param request body model.UpdatePostRequest true "Updated post data"
// @Success 200 {object} model.Post
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	post, err := h.postService.UpdatePost(postID, userID, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post (only by owner)
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	err := h.postService.DeletePost(postID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetFeed godoc
// @Summary Get user feed
// @Description Get posts from followed users
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.PostsResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/feed [get]
func (h *PostHandler) GetFeed(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.postService.GetFeed(userID, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserPosts godoc
// @Summary Get user posts
// @Description Get posts by a specific user
// @Tags posts
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.PostsResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/user/{user_id} [get]
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	viewerID := middleware.GetUserIDPtr(c)

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.postService.GetUserPosts(targetUserID, viewerID, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// LikePost godoc
// @Summary Like a post
// @Description Like a post
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/like [post]
func (h *PostHandler) LikePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	err := h.postService.LikePost(postID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post liked successfully",
		"liked":   true,
	})
}

// UnlikePost godoc
// @Summary Unlike a post
// @Description Unlike a post
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/unlike [post]
func (h *PostHandler) UnlikePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	err := h.postService.UnlikePost(postID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
		"liked":   false,
	})
}

// ToggleLike godoc
// @Summary Toggle post like
// @Description Like or unlike a post
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/toggle-like [post]
func (h *PostHandler) ToggleLike(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	liked, err := h.postService.ToggleLike(postID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	message := "Post unliked successfully"
	if liked {
		message = "Post liked successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"liked":   liked,
	})
}

// GetExplorePosts godoc
// @Summary Get explore posts
// @Description Get posts for exploration/discovery
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.PostsResponse
// @Security BearerAuth
// @Router /posts/explore [get]
func (h *PostHandler) GetExplorePosts(c *gin.Context) {
	userID := middleware.GetUserIDPtr(c)

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.postService.GetExplorePosts(userID, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPostStats godoc
// @Summary Get post statistics
// @Description Get statistics for a post
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/stats [get]
func (h *PostHandler) GetPostStats(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	stats, err := h.postService.GetPostStats(postID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// SearchPosts godoc
// @Summary Search posts
// @Description Search posts by content
// @Tags posts
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.PostsResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/search [get]
func (h *PostHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	userID := middleware.GetUserIDPtr(c)

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.postService.SearchPosts(query, userID, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
