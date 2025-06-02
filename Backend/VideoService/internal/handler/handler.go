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
    videos, err := service.GetAllVideosForUser(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) GetVideoForUser(c *gin.Context) {
    userID := c.Param("userID")
    videoID := c.Param("videoID")
    video, err := service.GetVideoForUser(userID, videoID)
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
    video, err := service.CreateVideoForUser(userID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, video)
}

func (h *VideoHandler) UpdateVideoForUser(c *gin.Context) {
    userID := c.Param("userID")
    videoID := c.Param("videoID")
    var req DTO.UpdateVideoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    video, err := service.UpdateVideoForUser(userID, videoID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, video)
}

func (h *VideoHandler) DeleteVideoForUser(c *gin.Context) {
    userID := c.Param("userID")
    videoID := c.Param("videoID")
    err := service.DeleteVideoForUser(userID, videoID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Video deleted"})
}