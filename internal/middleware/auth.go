package middleware

import (
    "net/http"
    "strings"
    "BAZ/internal/utils"
    "github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
    jwtUtil *utils.JWTUtil
}

func NewAuthMiddleware(jwtUtil *utils.JWTUtil) *AuthMiddleware {
    return &AuthMiddleware{jwtUtil: jwtUtil}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "missing authorization header"})
            c.Abort()
            return
        }
        
        // "Bearer <token>" extrahieren
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid authorization header format"})
            c.Abort()
            return
        }
        
        token := parts[1]
        claims, err := m.jwtUtil.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid or expired token"})
            c.Abort()
            return
        }
        
        // User-Informationen im Context speichern
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}

// Helper um UserID aus Context zu holen
func GetUserID(c *gin.Context) int64 {
    userID, exists := c.Get("userID")
    if !exists {
        return 0
    }
    return userID.(int64)
}