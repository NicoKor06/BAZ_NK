package main

import (
	"BAZ/internal/cache/redis"
	"BAZ/internal/db"
	"BAZ/internal/handler"
	"BAZ/internal/middleware"
	"BAZ/internal/oauth"
	"BAZ/internal/repository/postgres"
	"BAZ/internal/router"
	"BAZ/internal/usecase"
	"BAZ/internal/utils"
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// 1. CONFIGURATION & ENVIRONMENT
	envPath := filepath.Join(".", "env")

	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-me-in-production"
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable!")
	}
	jwtUtil := utils.NewJWTUtil(jwtSecret, 24)

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	// 2. INFRASTRUCTURE / EXTERNAL SERVICES (DB, Cache)
	// 2.1 Database
	conn, err := db.NewConnection(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(ctx)

	queries := db.NewQueriesFromConn(conn)

	// 2.2 Cache (Redis)
	redisCache, err := redis.NewRedisCache(redisHost, redisPort, "", 0)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisCache.Close()

	// 3. REPOSITORIES (Data Layer)
	userRepo := postgres.NewUserRepository(queries)
	blogRepo := postgres.NewBlogRepository(queries)
	commentRepo := postgres.NewCommentRepository(queries)

	// 4. USE CASES (Business Logic Layer)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil)
	userUsecase := usecase.NewUserUsecase(userRepo, blogRepo, commentRepo)
	blogUsecase := usecase.NewBlogUsecase(blogRepo, commentRepo, redisCache)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, blogRepo)
	oauthUsecase := usecase.NewOAuthUsecase(userRepo, jwtUtil)

	if clientID := os.Getenv("GITHUB_CLIENT_ID"); clientID != "" {
		githubProvider := oauth.NewGitHubProvider(
			clientID,
			os.Getenv("GITHUB_CLIENT_SECRET"),
			"http://localhost:8080/auth/github/callback",
		)
		oauthUsecase.RegisterProvider(githubProvider)
		log.Println("✅ GitHub OAuth2 enabled")
	}

	if clientID := os.Getenv("GOOGLE_CLIENT_ID"); clientID != "" {
		googleProvider := oauth.NewGoogleProvider(
			clientID,
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			"http://localhost:8080/auth/google/callback",
		)
		oauthUsecase.RegisterProvider(googleProvider)
		log.Println("Google OAuth2 enabled")
	}

	// 5. MIDDLEWARES & HANDLERS (Presentation Layer)
	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)
	cacheMiddleware := middleware.CacheMiddleware(redisCache, 5*time.Minute)
	rateLimiter := middleware.NewRateLimiter(redisCache, 100, time.Minute)

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	blogHandler := handler.NewBlogHandler(blogUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)
	oauthHandler := handler.NewOAuthHandler(oauthUsecase)

	// 6. ROUTER & SERVER START
	appRouter := router.NewRouter(
		authHandler,
		userHandler,
		blogHandler,
		commentHandler,
		oauthHandler,
		authMiddleware,
		cacheMiddleware,
		rateLimiter,
	)

	engine := appRouter.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
