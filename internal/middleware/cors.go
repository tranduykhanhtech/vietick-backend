package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures CORS settings for the API
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",    // React development server
			"http://localhost:3001",    // Alternative React port
			"https://vietick.com",      // Production frontend
			"https://www.vietick.com",  // Production frontend with www
			"https://app.vietick.com",  // App subdomain
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Total-Count",
			"X-Page",
			"X-Page-Size",
			"X-Has-More",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// DevelopmentCORSMiddleware provides more permissive CORS for development
func DevelopmentCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
