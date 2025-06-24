package service

import (
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"VideoService/internal/s3util"
	"log"
	"time"

	"github.com/google/uuid"
)

type VideoThumbInfo struct {
    VideoID        string `json:"videoID"`
    Title          string `json:"title"`
    TotalLikeCount int64  `json:"totalLikeCount"` 
    TotalViewCount uint64 `json:"totalViewCount"`
    ThumbnailURL   string `json:"thumbnailURL"`
}

func GetAllVideosForUser(userID uuid.UUID) ([]VideoThumbInfo, error) {
    log.Printf("[DEBUG] GetAllVideosForUser called for userID=%s", userID.String())
    var videos []entity.Video
    if err := database.DB.Where("user_id = ?", userID).Preload("Comments").Find(&videos).Error; err != nil {
        log.Printf("[ERROR] DB error in GetAllVideosForUser: %v", err)
        return nil, err
    }
    log.Printf("[DEBUG] Found %d videos for userID=%s", len(videos), userID.String())

    var results []VideoThumbInfo
    for _, video := range videos {
        log.Printf("[DEBUG] Processing videoID=%s", video.ID.String())
        thumbKey := userID.String() + "/" + video.ID.String() + "/thumbnail.jpeg"
        url, err := s3util.GeneratePresignedGetURL("toktikp2-video", thumbKey, 20*time.Minute)
        if err != nil {
            log.Printf("[WARN] Could not generate presigned URL for %s: %v", thumbKey, err)
            continue
        }
        log.Printf("[DEBUG] Generated thumbnail URL for videoID=%s: %s", video.ID.String(), url)
        results = append(results, VideoThumbInfo{
            VideoID:      video.ID.String(),
            ThumbnailURL: url,
            Title:        video.Title,
            TotalLikeCount: int64(video.TotalLikeCount),
            TotalViewCount: uint64(video.TotalViewCount),
        })
    }
    log.Printf("[DEBUG] Returning %d video thumbnails for userID=%s", len(results), userID.String())
    log.Printf("Returning videos: %+v", results)
    return results, nil
}