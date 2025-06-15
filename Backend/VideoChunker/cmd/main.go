package main

import (
    "fmt"
    "os"

	"VideoChuncker/internal/machineryUtil"
)

func main() {
    server, err := machineryUtil.CreateMachineryServer()
    if err != nil {
        fmt.Printf("Failed to create machinery server: %v\n", err)
        os.Exit(1)
    }

    // Use the StartWorker helper to register and launch the worker
    if err := machineryUtil.StartWorker(server); err != nil {
        fmt.Printf("Worker failed: %v\n", err)
        os.Exit(1)
    }
}