package main

import (
	"fmt"
	"log"
	"time"

	"vietick-backend/internal/config"
	"vietick-backend/internal/handler"
	"vietick-backend/internal/middleware"
	"vietick-backend/internal/repository"
	"vietick-backend/internal/service"
	"vietick-backend/pkg/database"
	"vietick-backend/pkg/email"
	"vietick-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

var version = "1.2.4" // Application version

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Loaded config: Server=%s:%s, DB=%s", cfg.Server.Host, cfg.Server.Port, cfg.Database.Name)

	// Initialize database
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Auto migrate database schema
	//  ###Production can not use auto migrate schema
	// if err := db.AutoMigrate(
	// 	&model.User{},
	// 	&model.Comment{},
	// 	&model.Post{},
	// 	&model.Follow{},
	// 	&model.IdentityVerification{},
	// 	&model.RefreshToken{},
	// ); err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiryHour,
		cfg.JWT.RefreshExpiryDay,
	)
	log.Println("JWT manager initialized")

	// Initialize email service
	emailService := email.NewEmailService(&cfg.Email)
	log.Println("Email service initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	followRepo := repository.NewFollowRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, authRepo, jwtManager, emailService)
	userService := service.NewUserService(userRepo, followRepo)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo)
	followService := service.NewFollowService(followRepo)
	verificationService := service.NewVerificationService(verificationRepo, userRepo, emailService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	followHandler := handler.NewFollowHandler(followService)
	verificationHandler := handler.NewVerificationHandler(verificationService)

	// Setup router
	router := setupRouter(cfg, authService, userService, authHandler, userHandler, postHandler, commentHandler, followHandler, verificationHandler)

	// Start cleanup routine for expired tokens
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := authService.CleanupExpiredTokens(); err != nil {
				log.Printf("Failed to cleanup expired tokens: %v", err)
			}
		}
	}()

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 VietTick Backend Server starting on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(
	cfg *config.Config,
	authService *service.AuthService,
	userService *service.UserService,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
	followHandler *handler.FollowHandler,
	verificationHandler *handler.VerificationHandler,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Global rate limiting
	router.Use(middleware.GlobalRateLimitMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "vietick-backend",
			"version":   version,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("")
		{
			// Auth routes with stricter rate limiting
			authGroup := public.Group("/auth")
			authGroup.Use(middleware.AuthRateLimitMiddleware())
			{
				authGroup.POST("/register", authHandler.Register)
				authGroup.POST("/login", authHandler.Login)
				authGroup.POST("/refresh", authHandler.RefreshToken)
				authGroup.POST("/verify-email", authHandler.VerifyEmail)
			}

			// Public user routes
			public.GET("/users/check-username", userHandler.CheckUsernameAvailability)
			public.GET("/users/check-email", userHandler.CheckEmailAvailability)
			public.GET("/verification/requirements", verificationHandler.GetVerificationRequirements)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		protected.Use(middleware.ApiRateLimitMiddleware())
		{
			// Auth routes (authenticated)
			authGroup := protected.Group("/auth")
			{
				authGroup.POST("/logout", authHandler.Logout)
				authGroup.POST("/logout-all", authHandler.LogoutAll)
				authGroup.POST("/resend-verification", authHandler.ResendEmailVerification)
				authGroup.POST("/change-password", authHandler.ChangePassword)
				authGroup.GET("/me", authHandler.GetProfile)
				authGroup.GET("/check", authHandler.CheckToken)
			}

			// User routes
			userGroup := protected.Group("/users")
			{
				userGroup.GET("/me", userHandler.GetCurrentProfile)
				userGroup.PUT("/me", userHandler.UpdateProfile)
				userGroup.PUT("/me/username", userHandler.UpdateUsername)
				userGroup.PUT("/me/email", userHandler.UpdateEmail)
				userGroup.GET("/recommended", userHandler.GetRecommendedUsers)
				userGroup.GET("/search", userHandler.SearchUsers)
				userGroup.GET("/:id", userHandler.GetProfile)
				userGroup.GET("/:id/stats", userHandler.GetUserStats)
				userGroup.GET("/username/:username", userHandler.GetProfileByUsername)

				// Follow routes
				userGroup.POST("/:id/follow", followHandler.Follow)
				userGroup.POST("/:id/unfollow", followHandler.Unfollow)
				userGroup.POST("/:id/toggle-follow", followHandler.ToggleFollow)
				userGroup.GET("/:id/follow-status", followHandler.GetFollowStatus)
				userGroup.GET("/:id/followers", followHandler.GetFollowers)
				userGroup.GET("/:id/following", followHandler.GetFollowing)
				userGroup.GET("/:id/follow-counts", followHandler.GetFollowCounts)
				userGroup.GET("/:id/mutual-follows", followHandler.GetMutualFollows)
				userGroup.GET("/:id/relationship", followHandler.GetFollowRelationship)
				userGroup.GET("/:id/follow-stats", followHandler.GetFollowStats)
			}

			// Follow routes (bulk operations)
			followGroup := protected.Group("/follows")
			{
				followGroup.POST("/bulk-follow", followHandler.BulkFollow)
				followGroup.POST("/bulk-unfollow", followHandler.BulkUnfollow)
			}

			// Post routes
			postGroup := protected.Group("/posts")
			{
				postGroup.POST("", postHandler.CreatePost)
				postGroup.GET("/feed", postHandler.GetFeed)
				postGroup.GET("/explore", postHandler.GetExplorePosts)
				postGroup.GET("/search", postHandler.SearchPosts)
				postGroup.GET("/:id", postHandler.GetPost)
				postGroup.PUT("/:id", postHandler.UpdatePost)
				postGroup.DELETE("/:id", postHandler.DeletePost)
				postGroup.GET("/:id/stats", postHandler.GetPostStats)
				postGroup.POST("/:id/like", postHandler.LikePost)
				postGroup.POST("/:id/unlike", postHandler.UnlikePost)
				postGroup.POST("/:id/toggle-like", postHandler.ToggleLike)
				postGroup.GET("/user/:user_id", postHandler.GetUserPosts)

				// Comment routes
				postGroup.POST("/:id/comments", commentHandler.CreateComment)
				postGroup.GET("/:id/comments", commentHandler.GetPostComments)
			}

			// Comment routes
			commentGroup := protected.Group("/comments")
			{
				commentGroup.GET("/:id", commentHandler.GetComment)
				commentGroup.PUT("/:id", commentHandler.UpdateComment)
				commentGroup.DELETE("/:id", commentHandler.DeleteComment)
				commentGroup.GET("/:id/stats", commentHandler.GetCommentStats)
				commentGroup.POST("/:id/like", commentHandler.LikeComment)
				commentGroup.POST("/:id/unlike", commentHandler.UnlikeComment)
				commentGroup.POST("/:id/toggle-like", commentHandler.ToggleLike)
			}

			// Verification routes
			verificationGroup := protected.Group("/verification")
			{
				verificationGroup.GET("/me", verificationHandler.GetUserVerification)
				verificationGroup.GET("/can-submit", verificationHandler.CanSubmitVerification)
				verificationGroup.GET("/verified-users", verificationHandler.GetVerifiedUsers)

				// Routes requiring email verification
				emailVerified := verificationGroup.Group("")
				emailVerified.Use(middleware.RequireEmailVerificationMiddleware(userService))
				{
					emailVerified.POST("/submit", verificationHandler.SubmitIdentityVerification)
				}

				// Admin routes
				adminRoutes := verificationGroup.Group("")
				adminRoutes.Use(middleware.AdminMiddleware(userService))
				{
					adminRoutes.GET("/pending", verificationHandler.GetPendingVerifications)
					adminRoutes.GET("/all", verificationHandler.GetAllVerifications)
					adminRoutes.GET("/stats", verificationHandler.GetVerificationStats)
					adminRoutes.GET("/:id", verificationHandler.GetVerification)
					adminRoutes.POST("/:id/review", verificationHandler.ReviewVerification)
					adminRoutes.DELETE("/:id", verificationHandler.DeleteVerification)
				}
			}
		}

		// Optional auth routes (can work with or without authentication)
		optionalAuth := v1.Group("")
		optionalAuth.Use(middleware.OptionalAuthMiddleware(authService))
		{
			// These routes can provide different data based on authentication status
			// For example, showing follow status if authenticated
		}
	}

	// 404 handler
	router.NoRoute(middleware.NotFoundMiddleware())

	return router
}
