package main

import (
	"context"
	"log"
	"os"

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

	// Datenbankverbindung mit pgx
	conn, err := db.NewConnection(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close(ctx)

	// Queries aus der Verbindung
	queries := db.NewQueriesFromConn(conn)

	// Repositories
	userRepo := postgres.NewUserRepository(queries)
	blogRepo := postgres.NewBlogRepository(queries)
	commentRepo := postgres.NewCommentRepository(queries)

	// JWT Utility
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-me-in-production"
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable!")
	}
	jwtUtil := utils.NewJWTUtil(jwtSecret, 24)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil)
	userUsecase := usecase.NewUserUsecase(userRepo, blogRepo, commentRepo)
	blogUsecase := usecase.NewBlogUsecase(blogRepo, commentRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, blogRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	blogHandler := handler.NewBlogHandler(blogUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)

	// Router
	appRouter := router.NewRouter(
		authHandler,
		userHandler,
		blogHandler,
		commentHandler,
		authMiddleware,
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
