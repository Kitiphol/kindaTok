package main

import (
	
    "Auth/internal/database"
    "Auth/internal/routes"
	"Auth/internal/entity"
	"log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {

    database.InitDB()
	
	seedAdmin()


    r := routes.Setup()
    r.Run() // default listens on :8080
}

func seedAdmin() {
    db := database.DB

    // Check if “admin” already exists
    var count int64
    db.Model(&entity.User{}).Where("username = ?", "admin").Count(&count)
    if count > 0 {
        log.Println("Admin user already exists, skipping seed.")
        return
    }

    // Hash the password “1”
    hashed, err := bcrypt.GenerateFromPassword([]byte("1"), bcrypt.DefaultCost)
    if err != nil {
        log.Fatalf("Failed to hash admin password: %v", err)
    }

    // Create the admin record
    admin := entity.User{
        ID:           uuid.New(),
        Username:     "admin",
        PasswordHash: string(hashed),
    }
    if err := db.Create(&admin).Error; err != nil {
        log.Fatalf("Failed to create admin user: %v", err)
    }
    log.Println("Seeded admin user (username=admin, password=1)")
}