package main

import (
	"context"
	"log"
	"os"
	"time"

	"BAZ/internal/cache"
	"BAZ/internal/db"
	"BAZ/internal/handler"
	"BAZ/internal/middleware"
	"BAZ/internal/repository/postgres"
	"BAZ/internal/router"
	"BAZ/internal/usecase"
	"BAZ/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ctx := context.Background()

	// ========== DATENBANK ==========
	conn, err := db.NewConnection(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(ctx)

	queries := db.NewQueriesFromConn(conn)

	// ========== REPOSITORIES ==========
	userRepo := postgres.NewUserRepository(queries)
	blogRepo := postgres.NewBlogRepository(queries)
	commentRepo := postgres.NewCommentRepository(queries)

	// ========== JWT ==========
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-me-in-production"
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable!")
	}
	jwtUtil := utils.NewJWTUtil(jwtSecret, 24)

	// ========== REDIS CACHE ==========
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisCache, err := cache.NewRedisCache(redisHost, redisPort, "", 0)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisCache.Close()

	// ========== USECASES (mit Cache) ==========
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil)
	userUsecase := usecase.NewUserUsecase(userRepo, blogRepo, commentRepo)
	blogUsecase := usecase.NewBlogUsecase(blogRepo, commentRepo, redisCache) // ← nur HIER!
	commentUsecase := usecase.NewCommentUsecase(commentRepo, blogRepo)

	// ========== HANDLER ==========
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	blogHandler := handler.NewBlogHandler(blogUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)

	// ========== MIDDLEWARE ==========
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)
	cacheMiddleware := middleware.CacheMiddleware(redisCache, 5*time.Minute)

	// ========== Redis ==========
	ratelimiter := middleware.NewRateLimiter(redisCache, 100, time.Minute)

	// ========== ROUTER ==========
	appRouter := router.NewRouter(
		authHandler,
		userHandler,
		blogHandler,
		commentHandler,
		authMiddleware,
		cacheMiddleware,
		ratelimiter,
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
