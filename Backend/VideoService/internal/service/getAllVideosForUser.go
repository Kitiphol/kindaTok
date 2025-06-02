package service

import (
    "VideoService/internal/entity"
    "VideoService/internal/database"
)

func GetAllVideosForUser(userID string) ([]entity.Video, error) {
    var videos []entity.Video
    err := database.DB.Where("user_id = ?", userID).Preload("Comments").Find(&videos).Error
    return videos, err
}