
package service

import (
	"errors"
    "VideoService/internal/entity"
    "VideoService/internal/database"
)

func GetVideoForUser(userID, videoID string) (entity.Video, error) {
    var video entity.Video
    err := database.DB.Where("user_id = ? AND id = ?", userID, videoID).Preload("Comments").First(&video).Error
    if err != nil {
        return video, errors.New("video not found")
    }
    return video, nil
}