package service

import (
    "VideoService/internal/entity"
    "VideoService/internal/database"
)

func GetAllVideos() ([]entity.Video, error) {
    var videos []entity.Video
    err := database.DB.Preload("User").Preload("Comments").Find(&videos).Error
    return videos, err
}