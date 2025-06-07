package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "VideoService/internal/service"
    "VideoService/internal/DTO"
)

type VideoHandler struct{}

func (h *VideoHandler) GetAllVideos(c *gin.Context) {
    videos, err := service.GetAllVideos()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) GetAllVideosForUser(c *gin.Context) {
    userID := c.Param("userID")

    //this presignURL should be Thumbnail URLs
    presignURL, err := service.GetAllVideosForUser(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, presignURL)
}

func (h *VideoHandler) GetVideoForUser(c *gin.Context) {
    userID := c.Param("userID")
    videoID := c.Param("videoID")
    presignURL, err := service.GetVideoForUser(userID, videoID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, presignURL)
}

func (h *VideoHandler) GetVideoForAllUser(c *gin.Context) {
    videoID := c.Param("videoID")
    video, err := service.GetVideoForAllUser(videoID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, video)
}

func (h *VideoHandler) CreateVideoForUser(c *gin.Context) {
    userID := c.Param("userID")
    var req DTO.CreateVideoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
   presignURL, err := service.CreateVideoForUser(userID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Video created successfully", "presignURL": presignURL})

}

func (h *VideoHandler) UpdateVideoForUser(c *gin.Context) {
    userID := c.Param("userID")
    videoID := c.Param("videoID")
    var req DTO.UpdateVideoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := service.UpdateVideoForUser(userID, videoID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Video updated successfully"} )
}

func (h *VideoHandler) DeleteVideoForUser(c *gin.Context) {
    userID := c.Param("userID")
    videoID := c.Param("videoID")
    err := service.DeleteVideoForUser(userID, videoID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    
    // Optionally, you can also delete the video file from S3 here
    
    c.JSON(http.StatusOK, gin.H{"message": "Video deleted"})
}