package service

import (
    "VideoService/internal/entity"
    "VideoService/internal/database"
    "github.com/google/uuid"

)

func GetAllVideosForUser(userID uuid.UUID) ([]entity.Video, error) {
    var videos []entity.Video
    if err := database.DB.Where("user_id = ?", userID).Preload("Comments").Find(&videos).Error; err != nil {
        return nil, err
    }
    return videos, nil
}