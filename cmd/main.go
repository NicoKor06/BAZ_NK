package main

import (
	"log"
	"os"

	"BAZ/internal/handler"
	"BAZ/internal/middleware"
	"BAZ/internal/router"
	"BAZ/internal/usecase"
	"BAZ/internal/utils"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
	}

	jwtUtil := utils.NewJWTUtil(jwtSecret, 24)

	authUsecase := &usecase.AuthUsecase{}
	userUsecase := &usecase.UserUsecase{}
	blogUsecase := &usecase.BlogUsecase{}
	commentUsecase := &usecase.CommentUsecase{}

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	blogHandler := handler.NewBlogHandler(blogUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)

	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)

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
