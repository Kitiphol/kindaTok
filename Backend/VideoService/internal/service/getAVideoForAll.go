
package service

import (
	"errors"
    "VideoService/internal/entity"
    "VideoService/internal/database"
	"VideoService/internal/s3util"
)


//
func GetVideoForAllUser(videoID string) (string, error) {
    var video entity.Video
    err := database.DB.Where("id = ?", videoID).Preload("Comments").First(&video).Error
    if err != nil {
        return "", errors.New("video not found")
    }
	video.ViewCount++
	err = database.DB.Save(&video).Error
	if err != nil {
		return "", errors.New("Failed to update view count for All Users")
	}

	presignedURL, err := s3util.GeneratePresignedGetURL("toktikp2-video", video.Filename, 90)
    return presignedURL, nil
}