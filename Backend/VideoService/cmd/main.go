package main

import (
    "VideoService/internal/database"
	"VideoService/internal/config"
	"VideoService/internal/routes"


)

func main() {
    database.InitDB()

    serverPort := config.Load().Port
    r := routes.Setup()
    r.Run(":" + serverPort)
}