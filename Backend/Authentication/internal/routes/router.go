package routes

import (
    "github.com/gin-gonic/gin"
    "Auth/internal/handler"
)

// Setup creates the Gin router and routes
type RegisterHandler = handler.RegisterHandler

func Setup() *gin.Engine {
    r := gin.Default()

    auth := r.Group("/api")
    reg := RegisterHandler{}
    auth.POST("/register", reg.Handle)
    auth.POST("/login", reg.Login)

    return r
}