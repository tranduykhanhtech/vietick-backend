package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details"`
}

// ValidationMiddleware handles validation errors
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were validation errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			if validationErr, ok := err.Err.(validator.ValidationErrors); ok {
				details := make(map[string]string)
				for _, fieldErr := range validationErr {
					details[fieldErr.Field()] = getValidationErrorMessage(fieldErr)
				}

				c.JSON(http.StatusBadRequest, ValidationErrorResponse{
					Error:   "Validation failed",
					Details: details,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		}
	}
}

// getValidationErrorMessage returns user-friendly error messages for validation errors
func getValidationErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + fieldErr.Param() + " characters long"
	case "max":
		return "Must be at most " + fieldErr.Param() + " characters long"
	case "oneof":
		return "Must be one of: " + fieldErr.Param()
	case "uuid":
		return "Must be a valid UUID"
	case "url":
		return "Must be a valid URL"
	default:
		return "Invalid value"
	}
}

// IDParamMiddleware validates ID parameters in URL
func IDParamMiddleware(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param(paramName)
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID parameter is required",
			})
			c.Abort()
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID parameter",
			})
			c.Abort()
			return
		}

		c.Set(paramName, id)
		c.Next()
	}
}

// BindJSONMiddleware validates and binds JSON request body
func BindJSONMiddleware(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			if validationErr, ok := err.(validator.ValidationErrors); ok {
				details := make(map[string]string)
				for _, fieldErr := range validationErr {
					details[fieldErr.Field()] = getValidationErrorMessage(fieldErr)
				}

				c.JSON(http.StatusBadRequest, ValidationErrorResponse{
					Error:   "Validation failed",
					Details: details,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON format",
			})
			c.Abort()
			return
		}

		c.Set("validated_body", obj)
		c.Next()
	}
}

// PaginationMiddleware validates pagination parameters
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("page_size", "20")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid page parameter",
			})
			c.Abort()
			return
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid page_size parameter (must be between 1 and 100)",
			})
			c.Abort()
			return
		}

		c.Set("page", page)
		c.Set("page_size", pageSize)
		c.Next()
	}
}

// ContentTypeMiddleware ensures the request has the correct content type
func ContentTypeMiddleware(expectedType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType != expectedType {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": "Content-Type must be " + expectedType,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequiredQueryMiddleware validates required query parameters
func RequiredQueryMiddleware(params ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, param := range params {
			value := c.Query(param)
			if value == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Required query parameter '" + param + "' is missing",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
