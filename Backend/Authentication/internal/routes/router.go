package routes

import (
	"Auth/internal/handler"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

// Setup creates the Gin router and routes
type RegisterHandler = handler.RegisterHandler

func Setup() *gin.Engine {
    r := gin.Default()

     // Configure CORS options
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // frontend origin
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))


    auth := r.Group("/api")
    reg := RegisterHandler{}
    auth.POST("/register", reg.Handle)
    auth.POST("/login", reg.Login)
    auth.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// auth.POST("/logout", reg.Logout)
    return r
}
