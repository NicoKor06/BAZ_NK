package router

import (
	"BAZ/internal/handler"
	"BAZ/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler     *handler.AuthHandler
	userHandler     *handler.UserHandler
	blogHandler     *handler.BlogHandler
	commentHandler  *handler.CommentHandler
	authMiddleware  *middleware.AuthMiddleware
	cacheMiddleware gin.HandlerFunc
	ratelimiter     *middleware.RateLimiter
}

// NewRouter erwartet jetzt 6 Parameter (vorher 5)
func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	blogHandler *handler.BlogHandler,
	commentHandler *handler.CommentHandler,
	authMiddleware *middleware.AuthMiddleware,
	cacheMiddleware gin.HandlerFunc,
	ratelimiter *middleware.RateLimiter,
) *Router {
	return &Router{
		authHandler:     authHandler,
		userHandler:     userHandler,
		blogHandler:     blogHandler,
		commentHandler:  commentHandler,
		authMiddleware:  authMiddleware,
		cacheMiddleware: cacheMiddleware,
		ratelimiter:     ratelimiter,
	}
}

func (r *Router) Setup() *gin.Engine {
	engine := gin.Default()
	engine.Use(r.ratelimiter.Middleware())
	// Öffentliche Routen (kein Token nötig)
	r.setupAuthRoutes(engine)
	r.setupPublicBlogRoutes(engine) // ← hier wird Cache angewendet
	r.setupPublicUserRoutes(engine)
	r.setupPublicCommentRoutes(engine)

	// Geschützte Routen (Token erforderlich)
	protected := engine.Group("")
	protected.Use(r.authMiddleware.Authenticate())
	{
		r.setupProtectedUserRoutes(protected)
		r.setupProtectedBlogRoutes(protected)
		r.setupProtectedCommentRoutes(protected)
	}

	return engine
}

// ========== ÖFFENTLICHE ROUTES ==========

func (r *Router) setupAuthRoutes(engine *gin.Engine) {
	auth := engine.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
	}
}

func (r *Router) setupPublicBlogRoutes(engine *gin.Engine) {
	blog := engine.Group("/blogs")
	{
		// ⭐ CACHE WIRD HIER ANGEWENDET (nur auf GET)
		blog.GET("", r.cacheMiddleware, r.blogHandler.GetAll)
		blog.GET("/:blogId", r.blogHandler.GetByID)
	}
}

func (r *Router) setupPublicUserRoutes(engine *gin.Engine) {
	engine.GET("/users/:userId", r.userHandler.GetPublicProfile)
}

func (r *Router) setupPublicCommentRoutes(engine *gin.Engine) {
	engine.GET("/blogs/:blogId/comments", r.commentHandler.GetByBlogID)
	engine.GET("/comments/:commentId", r.commentHandler.GetByID)
}

// ========== GESCHÜTZTE ROUTES (mit Token) ==========

func (r *Router) setupProtectedUserRoutes(group *gin.RouterGroup) {
	user := group.Group("/users")
	{
		user.GET("/me", r.userHandler.GetOwnProfile)
		user.DELETE("/me", r.userHandler.DeleteOwnAccount)
	}
}

func (r *Router) setupProtectedBlogRoutes(group *gin.RouterGroup) {
	blog := group.Group("/blogs")
	{
		blog.POST("", r.blogHandler.Create)
		blog.PUT("/:blogId", r.blogHandler.Update)
		blog.DELETE("/:blogId", r.blogHandler.Delete)
	}
}

func (r *Router) setupProtectedCommentRoutes(group *gin.RouterGroup) {
	group.POST("/blogs/:blogId/comments", r.commentHandler.Create)

	comment := group.Group("/comments")
	{
		comment.PUT("/:commentId", r.commentHandler.Update)
		comment.DELETE("/:commentId", r.commentHandler.Delete)
	}
}
