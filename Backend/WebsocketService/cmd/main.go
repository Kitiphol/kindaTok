package main

import (
    "WebsocketService/internal/database"
	"WebsocketService/internal/config"
	"WebsocketService/internal/routes"
    "WebsocketService/internal/machineryutil"

    "fmt"
    "os"
)


func main() {
    database.InitDB()

    serverPort := config.Load().Port
    r := routes.Setup()
    r.Run(":" + serverPort)

    _, err := machineryutil.CreateMachineryServer()
    if err != nil {
        fmt.Printf("Failed to create machinery server: %v\n", err)
        os.Exit(1)
    }

}