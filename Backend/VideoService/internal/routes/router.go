package routes

import (
    "github.com/gin-gonic/gin"
    "VideoService/internal/handler"
    "VideoService/internal/middleware"
)

// Setup creates the Gin router and routes
// type RegisterHandler = handler.RegisterHandler

func Setup() *gin.Engine {
    r := gin.Default()

	videoHandler := handler.VideoHandler{}

	r.GET("/videos", videoHandler.GetAllVideos)

    protected := r.Group("/", middleware.AuthMiddleware()) 
    {

        protected.GET("/videos/:videoID", videoHandler.GetVideoForAllUser)

        user := protected.Group("/users") 
        {

                // Get all videos for a user
            // by using :userID, you define it as a dynamic variable 
            // that can be parsed in gin in the handler functions
            user.GET("/:userID/videos", videoHandler.GetAllVideosForUser)

            // Get a specific video for a user
            user.GET("/:userID/videos/:videoID", videoHandler.GetVideoForUser)

            // Create a new video for a user
            user.POST("/:userID/videos", videoHandler.CreateVideoForUser)

            // (Update a specific video for a user
            user.PUT("/:userID/videos/:videoID", videoHandler.UpdateVideoForUser)

            // Delete a specific video for a user
            user.DELETE("/:userID/videos/:videoID", videoHandler.DeleteVideoForUser)

        }
    }
	


    
    return r
}