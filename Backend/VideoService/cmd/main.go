package main

import (
    // "log"
    // "github.com/google/uuid"
    "VideoService/internal/database"
    // "VideoService/internal/entity"
	"VideoService/internal/config"
	"VideoService/internal/routes"
)

func main() {
    database.InitDB()
    // seedVideos() 

    serverPort := config.Load().Port
    r := routes.Setup()
    r.Run(":" + serverPort)
}


// func seedVideos() {
//     db := database.DB

//     // Example: Replace with an actual user ID from your users table
//     userID := uuid.MustParse("ef4d776f-44d8-4ea9-97ed-33cd9f1f805e")

//     videos := []entity.Video{
//         {
//             ID:           uuid.New(),
//             Title:        "Sample Video 1",
//             UserID:       userID,
//             S3URL:        "https://example.com/video1.mp4",
//             ThumbnailURL: "https://example.com/thumb1.jpg",
//         },
//         {
//             ID:           uuid.New(),
//             Title:        "Sample Video 2",
//             UserID:       userID,
//             S3URL:        "https://example.com/video2.mp4",
//             ThumbnailURL: "https://example.com/thumb2.jpg",
//         },
//     }

//     for _, v := range videos {
//         if err := db.Create(&v).Error; err != nil {
//             log.Printf("Failed to create video: %v", err)
//         } else {
//             log.Printf("Seeded video: %s", v.Title)
//         }
//     }
// }