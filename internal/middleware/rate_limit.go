package middleware

import (
	"BAZ/internal/cache"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	cache  cache.Cache
	limit  int
	window time.Duration
}

func NewRateLimiter(cache cache.Cache, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		cache:  cache,
		limit:  limit,
		window: window,
	}
}

func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	fullkey := "ratelimit" + key

	countStr, err := r.cache.Get(ctx, fullkey)
	if err != nil {
		return false, err
	}

	if countStr == "" {
		r.cache.Set(ctx, fullkey, "1", r.window)
		return true, nil
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return false, err
	}

	if count >= r.limit {
		return false, nil
	}

	newCount := count + 1
	r.cache.Set(ctx, fullkey, strconv.Itoa(newCount), r.window)
	return true, nil
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Key = IP-Adresse (kann auch API-Key oder UserID sein)
		key := c.ClientIP()

		allowed, err := r.Allow(c.Request.Context(), key)
		if err != nil {
			c.Next() // Bei Fehler trotzdem durchlassen
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
