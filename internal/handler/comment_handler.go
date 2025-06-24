package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vietick-backend/internal/middleware"
	"vietick-backend/internal/model"
	"vietick-backend/internal/service"
	"vietick-backend/internal/utils"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment godoc
// @Summary Create a comment
// @Description Create a comment on a post
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param request body model.CreateCommentRequest true "Comment data"
// @Success 201 {object} model.Comment
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{post_id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	postID := c.Param("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	comment, err := h.commentService.CreateComment(userID, postID, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComment godoc
// @Summary Get a comment
// @Description Get a comment by ID
// @Tags comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} model.Comment
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	userID := middleware.GetUserIDPtr(c)

	comment, err := h.commentService.GetComment(commentID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetPostComments godoc
// @Summary Get post comments
// @Description Get comments for a post
// @Tags comments
// @Produce json
// @Param post_id path string true "Post ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.CommentsResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /posts/{post_id}/comments [get]
func (h *CommentHandler) GetPostComments(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	userID := middleware.GetUserIDPtr(c)

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.commentService.GetPostComments(postID, userID, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateComment godoc
// @Summary Update a comment
// @Description Update a comment (only by owner)
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param request body model.CreateCommentRequest true "Updated comment data"
// @Success 200 {object} model.Comment
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	comment, err := h.commentService.UpdateComment(commentID, userID, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment (only by owner)
// @Tags comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	err := h.commentService.DeleteComment(commentID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}

// LikeComment godoc
// @Summary Like a comment
// @Description Like a comment
// @Tags comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id}/like [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	err := h.commentService.LikeComment(commentID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment liked successfully",
		"liked":   true,
	})
}

// UnlikeComment godoc
// @Summary Unlike a comment
// @Description Unlike a comment
// @Tags comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id}/unlike [post]
func (h *CommentHandler) UnlikeComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	err := h.commentService.UnlikeComment(commentID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment unliked successfully",
		"liked":   false,
	})
}

// ToggleLike godoc
// @Summary Toggle comment like
// @Description Like or unlike a comment
// @Tags comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id}/toggle-like [post]
func (h *CommentHandler) ToggleLike(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	liked, err := h.commentService.ToggleLike(commentID, userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	message := "Comment unliked successfully"
	if liked {
		message = "Comment liked successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"liked":   liked,
	})
}

// GetCommentStats godoc
// @Summary Get comment statistics
// @Description Get statistics for a comment
// @Tags comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id}/stats [get]
func (h *CommentHandler) GetCommentStats(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	stats, err := h.commentService.GetCommentStats(commentID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}
