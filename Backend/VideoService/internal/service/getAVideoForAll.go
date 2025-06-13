
package service

import (
	"errors"
    "VideoService/internal/entity"
    "VideoService/internal/database"
	"VideoService/internal/s3util"
	"github.com/google/uuid"
)

func GetVideoForAllUser(videoID uuid.UUID) (string, error) {
    var video entity.Video
    if err := database.DB.Where("id = ?", videoID).Preload("Comments").First(&video).Error; err != nil {
        return "", errors.New("video not found")
    }
    video.ViewCount++
    if err := database.DB.Save(&video).Error; err != nil {
        return "", errors.New("failed to update view count")
    }
    return s3util.GeneratePresignedGetURL("toktikp2-video", video.Filename, 90)
}