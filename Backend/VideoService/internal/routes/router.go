    package routes

    import (
        "VideoService/internal/handler"
        "VideoService/internal/middleware"
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

        videoHandler := handler.VideoHandler{}

        // Public: list all videos
        r.GET("/api/videos", videoHandler.GetAllVideos)

        // Protected routes require valid JWT
        protected := r.Group("/api", middleware.AuthMiddleware())
        {
            // List videos belonging to the authenticated user
            protected.GET("/videos/user", videoHandler.GetAllVideosForUser)

            // Get a single video (must belong to the authenticated user)
            protected.GET("/videos/:videoID", videoHandler.GetVideoForUser)

            // Create a new video record & return presigned URL
            protected.POST("/videos", videoHandler.CreateVideoForUser)

            protected.GET("/videos/check", videoHandler.CheckUpload)

            // Update metadata for a user's video
            protected.PUT("/videos/:videoID", videoHandler.UpdateVideoForUser)

            // Delete a user's video
            protected.DELETE("/videos/:videoID", videoHandler.DeleteVideoForUser)
        }

        return r
    }
