package service

import (
	"VideoService/internal/DTO"
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"VideoService/internal/s3util"
	"fmt"

	"github.com/google/uuid"
)


func CreateVideoForUser(userID uuid.UUID, req DTO.CreateVideoRequest) (string, error) {
    videoID := uuid.New()

    // uuid_{filename}
    uniqueFilename := fmt.Sprintf("videos/%s_%s", videoID.String(), req.Filename)

    video := entity.Video{
        ID:           videoID,
        Title:        req.Title,
        Filename:     uniqueFilename,
        UserID:       userID,
    }
    if err := database.DB.Create(&video).Error; err != nil {
        return "", err
    }

    presignURL, err := s3util.GeneratePresignedPutURL("toktikp2-video", uniqueFilename, 300)
    if err != nil {
        return "", err
    }
    return presignURL, nil
}