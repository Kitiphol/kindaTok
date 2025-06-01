package main

import (
    
	"UserService/internal/database"
	"UserService/internal/routes"
)


func main() {

	// Initialize database connection
    database.InitDB()
	


    r := routes.Setup()
    r.Run(":8081")
}