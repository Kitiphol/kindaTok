package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "VideoService/internal/DTO"
    "VideoService/internal/middleware"
    "VideoService/internal/service"
)

type VideoHandler struct{}

// GetAllVideos returns all videos in the system (public)
func (h *VideoHandler) GetAllVideos(c *gin.Context) {
    videos, err := service.GetAllVideos()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

// GetAllVideosForUser returns videos for the authenticated user
func (h *VideoHandler) GetAllVideosForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videos, err := service.GetAllVideosForUser(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

// GetVideoForUser returns a presigned GET URL for a user's video
func (h *VideoHandler) GetVideoForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videoIDStr := c.Param("videoID")
    videoID, parseErr := uuid.Parse(videoIDStr)
    if parseErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }

    presignURL, err := service.GetVideoForUser(userID, videoID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"presignedURL": presignURL})
}

// CreateVideoForUser issues a presigned PUT URL for a new user video
func (h *VideoHandler) CreateVideoForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var req DTO.CreateVideoRequest
    if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
        return
    }

    presignURL, err := service.CreateVideoForUser(userID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "video record created", "presignedURL": presignURL})
}

// UpdateVideoForUser updates metadata for an existing video
func (h *VideoHandler) UpdateVideoForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videoIDStr := c.Param("videoID")
    videoID, parseErr := uuid.Parse(videoIDStr)
    if parseErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }

    var req DTO.UpdateVideoRequest
    if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
        return
    }

    if err := service.UpdateVideoForUser(userID, videoID, req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "video updated"})
}

// DeleteVideoForUser deletes a user's video
func (h *VideoHandler) DeleteVideoForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videoIDStr := c.Param("videoID")
    videoID, parseErr := uuid.Parse(videoIDStr)
    if parseErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }

    if err := service.DeleteVideoForUser(userID, videoID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "video deleted"})
}
