package service

import (
	"VideoService/internal/DTO"
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"VideoService/internal/s3util"
	"fmt"

	"github.com/google/uuid"
)


func CreateVideoForUser(userID string, req DTO.CreateVideoRequest) (string, error) {
	videoID := uuid.New()
    uniqueFilename := fmt.Sprintf("videos/%s_%s", videoID.String(), req.Filename)
	video := entity.Video{
        ID:           videoID,
        Title:        req.Title,
        Filename:     uniqueFilename,
        UserID:       uuid.MustParse(userID),
        // S3URL:        req.S3URL,
        // ThumbnailURL: req.ThumbnailURL,
    }
    err := database.DB.Create(&video).Error

	presignURL, err := s3util.GeneratePresignedPutURL("toktikp2-video", uniqueFilename, 300)
    return presignURL, err
}