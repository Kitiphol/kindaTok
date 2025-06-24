package routes

import (
    "WebsocketService/internal/service"
    "WebsocketService/internal/middleware"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

// Setup creates the Gin router and routes
func Setup() *gin.Engine {
    r := gin.Default()

    // Configure CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Protected routes require valid JWT
    protected := r.Group("/api", middleware.AuthMiddleware())
    {
        // update a comment on a video
        protected.POST("comments/videos/:videoID", service.AddOrUpdateComment)

        // update likes on a video
        protected.POST("likes/videos/:videoID", service.AddOrUpdateLike)

        // update a view count for a video
        protected.POST("view/videos/:videoID", service.AddOrUpdateView)

        // get the total likes for a video
        protected.GET("likes/videos/:videoID", service.GetTotalLikes)

        // get the total views for a video
        protected.GET("views/videos/:videoID", service.GetTotalViews)

        // list all comments for a video
        protected.GET("comments/videos/:videoID", service.ListComments)


        // Get a presigned GET URL for a user's video
        protected.DELETE("likes/videos/:videoID", service.DeleteLike)

        // Delete a comment for a video UserSpecified
        protected.DELETE("comments/videos/:videoID/:commentID", service.DeleteComment)
    }

    return r
}