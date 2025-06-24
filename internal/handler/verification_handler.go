package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vietick-backend/internal/middleware"
	"vietick-backend/internal/model"
	"vietick-backend/internal/service"
	"vietick-backend/internal/utils"
)

type VerificationHandler struct {
	verificationService *service.VerificationService
}

func NewVerificationHandler(verificationService *service.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
	}
}

// SubmitIdentityVerification godoc
// @Summary Submit identity verification
// @Description Submit identity verification for blue tick
// @Tags verification
// @Accept json
// @Produce json
// @Param request body model.SubmitIdentityVerificationRequest true "Identity verification data"
// @Success 201 {object} model.IdentityVerification
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/submit [post]
func (h *VerificationHandler) SubmitIdentityVerification(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.SubmitIdentityVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	verification, err := h.verificationService.SubmitIdentityVerification(userID, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, verification)
}

// GetUserVerification godoc
// @Summary Get user verification status
// @Description Get the current user's verification status
// @Tags verification
// @Produce json
// @Success 200 {object} model.IdentityVerification
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/me [get]
func (h *VerificationHandler) GetUserVerification(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	verification, err := h.verificationService.GetUserVerification(userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, verification)
}

// GetVerification godoc
// @Summary Get verification by ID
// @Description Get verification details by ID (admin only)
// @Tags verification
// @Produce json
// @Param id path string true "Verification ID"
// @Success 200 {object} model.IdentityVerification
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/{id} [get]
func (h *VerificationHandler) GetVerification(c *gin.Context) {
	verificationID := c.Param("id")
	if verificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification ID"})
		return
	}

	verification, err := h.verificationService.GetVerification(verificationID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, verification)
}

// GetPendingVerifications godoc
// @Summary Get pending verifications
// @Description Get all pending verification requests (admin only)
// @Tags verification
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.IdentityVerificationsResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/pending [get]
func (h *VerificationHandler) GetPendingVerifications(c *gin.Context) {
	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.verificationService.GetPendingVerifications(&pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllVerifications godoc
// @Summary Get all verifications
// @Description Get all verification requests with optional status filter (admin only)
// @Tags verification
// @Produce json
// @Param status query string false "Status filter" Enums(pending,approved,rejected)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.IdentityVerificationsResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/all [get]
func (h *VerificationHandler) GetAllVerifications(c *gin.Context) {
	statusStr := c.Query("status")
	var status *model.IdentityVerificationStatus
	if statusStr != "" {
		s := model.IdentityVerificationStatus(statusStr)
		status = &s
	}

	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.verificationService.GetAllVerifications(status, &pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ReviewVerification godoc
// @Summary Review verification request
// @Description Approve or reject a verification request (admin only)
// @Tags verification
// @Accept json
// @Produce json
// @Param id path string true "Verification ID"
// @Param request body model.ReviewIdentityVerificationRequest true "Review data"
// @Success 200 {object} model.IdentityVerification
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/{id}/review [post]
func (h *VerificationHandler) ReviewVerification(c *gin.Context) {
	reviewedBy, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	verificationID := c.Param("id")
	if verificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification ID"})
		return
	}

	var req model.ReviewIdentityVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, err)
		return
	}

	verification, err := h.verificationService.ReviewVerification(verificationID, reviewedBy, &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, verification)
}

// GetVerificationStats godoc
// @Summary Get verification statistics
// @Description Get verification statistics (admin only)
// @Tags verification
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/stats [get]
func (h *VerificationHandler) GetVerificationStats(c *gin.Context) {
	stats, err := h.verificationService.GetVerificationStats()
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// CanSubmitVerification godoc
// @Summary Check if user can submit verification
// @Description Check if the current user can submit a verification request
// @Tags verification
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/can-submit [get]
func (h *VerificationHandler) CanSubmitVerification(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	canSubmit, reason, err := h.verificationService.CanSubmitVerification(userID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"can_submit": canSubmit,
		"reason":     reason,
	})
}

// DeleteVerification godoc
// @Summary Delete verification request
// @Description Delete a verification request (admin only)
// @Tags verification
// @Produce json
// @Param id path string true "Verification ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Security BearerAuth
// @Router /verification/{id} [delete]
func (h *VerificationHandler) DeleteVerification(c *gin.Context) {
	verificationID := c.Param("id")
	if verificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification ID"})
		return
	}

	err := h.verificationService.DeleteVerification(verificationID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification deleted successfully",
	})
}

// GetVerifiedUsers godoc
// @Summary Get verified users
// @Description Get a list of verified users
// @Tags verification
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} model.FollowersResponse
// @Security BearerAuth
// @Router /verification/verified-users [get]
func (h *VerificationHandler) GetVerifiedUsers(c *gin.Context) {
	var pagination utils.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.HandleError(c, err)
		return
	}

	response, err := h.verificationService.GetVerifiedUsers(&pagination)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetVerificationRequirements godoc
// @Summary Get verification requirements
// @Description Get the requirements for identity verification
// @Tags verification
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /verification/requirements [get]
func (h *VerificationHandler) GetVerificationRequirements(c *gin.Context) {
	requirements := map[string]interface{}{
		"required_documents": []string{
			"Government-issued ID (front)",
			"Government-issued ID (back) - if applicable",
			"Selfie with ID",
		},
		"accepted_id_types": []string{
			"national_id",
			"passport",
			"driver_license",
		},
		"image_requirements": map[string]string{
			"format":     "JPEG, PNG",
			"max_size":   "5MB",
			"resolution": "Minimum 300x300 pixels",
			"quality":    "Clear and readable",
		},
		"processing_time": "3-5 business days",
		"prerequisites": []string{
			"Verified email address",
			"Complete profile information",
		},
		"note": "All information must match exactly between documents and profile",
	}

	c.JSON(http.StatusOK, requirements)
}
