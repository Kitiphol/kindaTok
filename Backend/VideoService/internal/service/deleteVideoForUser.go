package service

import (
	"VideoService/internal/database"
	"VideoService/internal/entity"
)


func DeleteVideoForUser(userID, videoID string) error {
    return database.DB.Where("user_id = ? AND id = ?", userID, videoID).Delete(&entity.Video{}).Error
}