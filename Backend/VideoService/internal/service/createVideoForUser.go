package service

import (
    "VideoService/internal/entity"
    "VideoService/internal/DTO"
    "VideoService/internal/database"
    "github.com/google/uuid"
)


func CreateVideoForUser(userID string, req DTO.CreateVideoRequest) (entity.Video, error) {
    video := entity.Video{
        ID:           uuid.New(),
        Title:        req.Title,
        // Filename:     req.Filename,
        UserID:       uuid.MustParse(userID),
        S3URL:        req.S3URL,
        ThumbnailURL: req.ThumbnailURL,
    }
    err := database.DB.Create(&video).Error
    return video, err
}