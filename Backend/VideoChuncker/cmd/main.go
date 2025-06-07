package main

import (
    "fmt"
    "os"

    "VideoChuncker/internal/service"
)

func main() {
    bucket := "toktikp2-video"                  
    objectKey := "stopwatch-1minute.mp4"          

    fmt.Printf("Starting chunk and upload for %s/%s...\n", bucket, objectKey)
    err := service.ChunkAndUpload(bucket, objectKey)
    if err != nil {
        fmt.Printf("Chunk and upload failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Chunk and upload succeeded!")
}