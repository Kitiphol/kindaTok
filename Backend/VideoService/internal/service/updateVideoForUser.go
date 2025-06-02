package service

import (
    "errors"
    "VideoService/internal/entity"
    "VideoService/internal/DTO"
    "VideoService/internal/database"
)

func UpdateVideoForUser(userID, videoID string, req DTO.UpdateVideoRequest) (entity.Video, error) {
    var video entity.Video
    err := database.DB.Where("user_id = ? AND id = ?", userID, videoID).First(&video).Error
    if err != nil {
        return video, errors.New("video not found")
    }
    if req.Title != "" {
        video.Title = req.Title
    }
    // if req.Filename != "" {
    //     video.Filename = req.Filename
    // }
    if req.S3URL != "" {
        video.S3URL = req.S3URL
    }
    if req.ThumbnailURL != "" {
        video.ThumbnailURL = req.ThumbnailURL
    }
    err = database.DB.Save(&video).Error
    return video, err
}