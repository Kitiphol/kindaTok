// internal/middleware/auth.go
package middleware

import (
    "errors"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
    "Auth/internal/config"
)

// Context key for userID
const ContextUserID = "userID"

// AuthMiddleware validates JWT, extracts userID, and sets it in Gin context
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(401, gin.H{"error": "Missing or malformed Authorization header"})
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        cfg := config.Load()
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(cfg.JWTSecret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["sub"] == nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
            return
        }

        userIDStr, ok := claims["sub"].(string)
        if !ok {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid userID claim"})
            return
        }
        userID, err := uuid.Parse(userIDStr)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Malformed userID in token"})
            return
        }

        // set userID into Gin context
        c.Set(ContextUserID, userID)
        c.Next()
    }
}

// GetUserIDFromContext retrieves the userID set by AuthMiddleware
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
    val, exists := c.Get(ContextUserID)
    if !exists {
        return uuid.Nil, errors.New("userID not found in context")
    }
    userID, ok := val.(uuid.UUID)
    if !ok {
        return uuid.Nil, errors.New("invalid userID type in context")
    }
    return userID, nil
}