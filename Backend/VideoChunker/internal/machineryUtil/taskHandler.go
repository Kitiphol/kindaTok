package machineryUtil

import (
    "context"
    "fmt"
    "github.com/RichardKnop/machinery/v2"
)

// Example task function
func ChunkVideo(ctx context.Context, input string) error {
    fmt.Println("Chunking video with input:", input)
    // Do your chunking logic here
    return nil
}

// StartWorker registers tasks and starts the worker
func StartWorker(server *machinery.Server) error {
    tasks := map[string]interface{}{
        "ChunkVideo": ChunkVideo,
    }
    if err := server.RegisterTasks(tasks); err != nil {
        return err
    }
    worker := server.NewWorker("video-chuncker-worker", 1)
    return worker.Launch()
}