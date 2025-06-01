package handler

import (
	"github.com/gin-gonic/gin"
	"UserService/internal/service"
	"UserService/internal/DTO"
)


// RegisterHandler handles POST /updateUser
type RegisterHandler struct{}
func (h *RegisterHandler) Handle(c *gin.Context) {
	var req DTO.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	msg, err := service.UpdateUser(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": msg})
}