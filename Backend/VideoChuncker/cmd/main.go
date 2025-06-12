package main

import (
    "fmt"
    "os"

    "VideoChuncker/internal/service"
	"VideoChuncker/internal/machineryUtil"
)

func main() {
    bucket := "toktikp2-video"                  
    objectKey := "stopwatch-1minute.mp4"          


	server, err := machineryutil.CreateMachineryServer()
    if err != nil {
        fmt.Printf("Failed to create machinery server: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Starting chunk and upload for %s/%s...\n", bucket, objectKey)
    err2 := service.ChunkAndUpload(server, bucket, objectKey)
    if err2 != nil {
        fmt.Printf("Chunk and upload failed: %v\n", err2)
        os.Exit(1)
    }
    fmt.Println("Chunk and upload succeeded!")
}