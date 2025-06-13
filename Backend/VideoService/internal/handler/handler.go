package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "VideoService/internal/service"
    "VideoService/internal/DTO"
    "VideoService/internal/middleware"
)

type VideoHandler struct{}

// GetAllVideos is public, no auth required
func (h *VideoHandler) GetAllVideos(c *gin.Context) {
    videos, err := service.GetAllVideos()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

// Helper to extract userID from context and abort on error
func getUserID(c *gin.Context) (uuid.UUID, bool) {
    id, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return uuid.Nil, false
    }
    return id, true
}

func (h *VideoHandler) GetAllVideosForUser(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }

    videos, err := service.GetAllVideosForUser(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) GetVideoForUser(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }

    presignURL, err := service.GetVideoForUser(userID, videoID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, presignURL)
}

func (h *VideoHandler) GetVideoForAllUser(c *gin.Context) {
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }

    videoURL, err := service.GetVideoForAllUser(videoID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videoURL)
}

func (h *VideoHandler) CreateVideoForUser(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }
    var req DTO.CreateVideoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    presignedURL, err := service.CreateVideoForUser(userID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Video created successfully", "presignedURL": presignedURL})
}

func (h *VideoHandler) UpdateVideoForUser(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }
    var req DTO.UpdateVideoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := service.UpdateVideoForUser(userID, videoID, req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Video updated successfully"})
}

func (h *VideoHandler) DeleteVideoForUser(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }
    if err := service.DeleteVideoForUser(userID, videoID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Video deleted"})
}