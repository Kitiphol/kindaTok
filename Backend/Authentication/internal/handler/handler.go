package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "Auth/internal/service"
	"Auth/internal/DTO"
)

// RegisterHandler handles POST /register
type RegisterHandler struct{}

func (h *RegisterHandler) Handle(c *gin.Context) {
    var req DTO.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    token, err := service.Register(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}

// LoginHandler handles POST /login
func (h *RegisterHandler) Login(c *gin.Context) {
    var req DTO.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    token, err := service.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}

// // POST /api/logout
// func (h *RegisterHandler) Logout(c *gin.Context) {
//     // Extract JWT from header
//     token := extractTokenFromHeader(c)
//     // Add token to blacklist (e.g., Redis)
//     blacklistToken(token)
//     c.JSON(200, gin.H{"message": "Logged out"})
// }