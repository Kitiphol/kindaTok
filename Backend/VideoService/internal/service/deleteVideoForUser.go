package service

import (
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"github.com/google/uuid"
)


func DeleteVideoForUser(userID, videoID uuid.UUID) error {
    return database.DB.Where("user_id = ? AND id = ?", userID, videoID).Delete(&entity.Video{}).Error
}