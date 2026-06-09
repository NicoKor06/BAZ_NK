package middleware

import (
	"bytes"
	"net/http"
	"time"

	"BAZ/internal/cache"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CacheMiddleware(cache cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		key := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			key = key + "?" + c.Request.URL.RawQuery
		}

		cached, err := cache.Get(c.Request.Context(), key)
		if err == nil && cached != "" {
			c.Data(http.StatusOK, "application/json", []byte(cached))
			c.Abort()
			return
		}

		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w
		c.Next()

		if c.Writer.Status() == http.StatusOK {
			cache.Set(c.Request.Context(), key, w.body.String(), ttl)
		}
	}
}
