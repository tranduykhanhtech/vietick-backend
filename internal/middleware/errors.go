package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error     string      `json:"error"`
	Message   string      `json:"message,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorHandlingMiddleware provides centralized error handling
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				fmt.Printf("Panic recovered: %v\n", err)
				fmt.Printf("Stack trace: %s\n", debug.Stack())

				// Get request ID if available
				requestID, _ := c.Get("request_id")
				requestIDStr := ""
				if requestID != nil {
					requestIDStr = requestID.(string)
				}

				// Return internal server error
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:     "Internal Server Error",
					Message:   "An unexpected error occurred",
					RequestID: requestIDStr,
				})
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors that were added to the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			HandleError(c, err.Err)
		}
	}
}

// HandleError provides consistent error handling
func HandleError(c *gin.Context, err error) {
	if c.IsAborted() {
		return // Response already sent
	}

	requestID, _ := c.Get("request_id")
	requestIDStr := ""
	if requestID != nil {
		requestIDStr = requestID.(string)
	}

	errorMsg := err.Error()
	statusCode := http.StatusInternalServerError
	message := ""

	// Determine status code based on error message
	switch {
	case strings.Contains(errorMsg, "not found"):
		statusCode = http.StatusNotFound
		message = "Resource not found"
	case strings.Contains(errorMsg, "unauthorized") || strings.Contains(errorMsg, "invalid token"):
		statusCode = http.StatusUnauthorized
		message = "Authentication required"
	case strings.Contains(errorMsg, "forbidden") || strings.Contains(errorMsg, "access denied"):
		statusCode = http.StatusForbidden
		message = "Access denied"
	case strings.Contains(errorMsg, "already exists") || strings.Contains(errorMsg, "duplicate"):
		statusCode = http.StatusConflict
		message = "Resource already exists"
	case strings.Contains(errorMsg, "validation") || strings.Contains(errorMsg, "invalid"):
		statusCode = http.StatusBadRequest
		message = "Invalid input"
	case strings.Contains(errorMsg, "rate limit"):
		statusCode = http.StatusTooManyRequests
		message = "Rate limit exceeded"
	default:
		message = "An error occurred"
	}

	c.JSON(statusCode, ErrorResponse{
		Error:     errorMsg,
		Message:   message,
		RequestID: requestIDStr,
	})
}

// NotFoundMiddleware handles 404 errors
func NotFoundMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: "The requested resource was not found",
		})
	}
}

// MethodNotAllowedMiddleware handles 405 errors
func MethodNotAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			Error:   "Method Not Allowed",
			Message: "The requested method is not allowed for this resource",
		})
	}
}

// TimeoutMiddleware handles request timeouts
func TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set timeout context if needed
		// This is a placeholder - in a real app you might want to implement actual timeout handling
		c.Next()
	}
}

// DatabaseErrorMiddleware handles database-specific errors
func DatabaseErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			errorMsg := err.Error()

			requestID, _ := c.Get("request_id")
			requestIDStr := ""
			if requestID != nil {
				requestIDStr = requestID.(string)
			}

			// Handle database-specific errors
			switch {
			case strings.Contains(errorMsg, "connection refused"):
				c.JSON(http.StatusServiceUnavailable, ErrorResponse{
					Error:     "Service Unavailable",
					Message:   "Database connection failed",
					RequestID: requestIDStr,
				})
			case strings.Contains(errorMsg, "timeout"):
				c.JSON(http.StatusRequestTimeout, ErrorResponse{
					Error:     "Request Timeout",
					Message:   "Database operation timed out",
					RequestID: requestIDStr,
				})
			case strings.Contains(errorMsg, "foreign key constraint"):
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:     "Constraint Violation",
					Message:   "Referenced resource does not exist",
					RequestID: requestIDStr,
				})
			}
		}
	}
}
