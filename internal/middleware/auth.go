package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"vietick-backend/internal/service"
	"vietick-backend/pkg/jwt"
)

// AuthMiddleware validates JWT tokens and sets user information in context
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must start with Bearer",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func OptionalAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.Next()
			return
		}

		// Validate token
		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			// Invalid token, but we don't fail the request
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireEmailVerificationMiddleware checks if user's email is verified
func RequireEmailVerificationMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		user, err := userService.GetUserByID(userID.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		if !user.IsEmailVerified {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Email verification required",
				"message": "Please verify your email address to access this feature",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireVerificationMiddleware checks if user is verified (has blue tick)
func RequireVerificationMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		user, err := userService.GetUserByID(userID.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account verification required",
				"message": "This feature requires a verified account",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware checks if user has admin privileges
// Note: This is a placeholder implementation. In a real app, you'd have a proper role system
func AdminMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		user, err := userService.GetUserByID(userID.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// Simple admin check - in a real app, you'd have a proper role system
		// For now, let's say a specific UUID is admin (replace with real admin UUID)
		const adminUUID = "00000000-0000-0000-0000-000000000001" // Example admin UUID
		if user.ID != adminUUID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID is a helper function to get user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserIDPtr is a helper function to get user ID pointer from context
func GetUserIDPtr(c *gin.Context) *string {
	userID, exists := GetUserID(c)
	if !exists {
		return nil
	}
	return &userID
}

// GetClaims is a helper function to get JWT claims from context
func GetClaims(c *gin.Context) (*jwt.Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	return claims.(*jwt.Claims), true
}
