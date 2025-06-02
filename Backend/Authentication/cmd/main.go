package main

import (
	
    "Auth/internal/database"
    "Auth/internal/routes"
	"Auth/internal/entity"
	"log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// psql -U dev -d dbp2 -c "SELECT * FROM users;"
// psql -U dev -d dbp2 -c "SELECT * FROM videos;"
func main() {

    database.InitDB()
	
	adminID := seedAdmin()
	seedVideos(adminID)


    r := routes.Setup()
    r.Run() // default listens on :8080
}


func seedAdmin() uuid.UUID {
    db := database.DB

    var admin entity.User
    // Try to find existing admin
    if err := db.Where("username = ?", "admin").First(&admin).Error; err == nil {
        log.Println("Admin user already exists, skipping seed.")
        return admin.ID
    }

    // Hash the password “1”
    hashed, err := bcrypt.GenerateFromPassword([]byte("1"), bcrypt.DefaultCost)
    if err != nil {
        log.Fatalf("Failed to hash admin password: %v", err)
    }

    // Create the admin record
    admin = entity.User{
        ID:           uuid.New(),
        Username:     "admin",
        PasswordHash: string(hashed),
    }
    if err := db.Create(&admin).Error; err != nil {
        log.Fatalf("Failed to create admin user: %v", err)
    }
    log.Println("Seeded admin user (username=admin, password=1)")
    return admin.ID
}

func seedVideos(userID uuid.UUID) {
    db := database.DB

    videos := []entity.Video{
        {
            ID:           uuid.New(),
            Title:        "Sample Video 1",
            UserID:       userID,
            S3URL:        "https://example.com/video1.mp4",
            ThumbnailURL: "https://example.com/thumb1.jpg",
        },
        {
            ID:           uuid.New(),
            Title:        "Sample Video 2",
            UserID:       userID,
            S3URL:        "https://example.com/video2.mp4",
            ThumbnailURL: "https://example.com/thumb2.jpg",
        },
    }

    for _, v := range videos {
        if err := db.Create(&v).Error; err != nil {
            log.Printf("Failed to create video: %v", err)
        } else {
            log.Printf("Seeded video: %s", v.Title)
        }
    }
}

// func seedAdmin() {
//     db := database.DB

//     // Check if “admin” already exists
//     var count int64
//     db.Model(&entity.User{}).Where("username = ?", "admin").Count(&count)
//     if count > 0 {
//         log.Println("Admin user already exists, skipping seed.")
//         return
//     }

//     // Hash the password “1”
//     hashed, err := bcrypt.GenerateFromPassword([]byte("1"), bcrypt.DefaultCost)
//     if err != nil {
//         log.Fatalf("Failed to hash admin password: %v", err)
//     }

//     // Create the admin record
//     admin := entity.User{
//         ID:           uuid.New(),
//         Username:     "admin",
//         PasswordHash: string(hashed),
//     }
//     if err := db.Create(&admin).Error; err != nil {
//         log.Fatalf("Failed to create admin user: %v", err)
//     }
//     log.Println("Seeded admin user (username=admin, password=1)")
// }


// func seedVideos() {
//     db := database.DB

//     // Replace with an actual user ID from your users table
//     userID := uuid.MustParse("PUT-A-REAL-USER-ID-HERE")

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