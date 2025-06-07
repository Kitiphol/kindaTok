package service

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "VideoChuncker/internal/s3util"
)


func ChunkAndUpload(bucket, objectKey string) error {

    // 1. Download the video from S3 to a local temp file
    localInput := filepath.Join(os.TempDir(), filepath.Base(objectKey))
    if err := s3util.DownloadFromR2(bucket, objectKey, localInput); err != nil {
        return fmt.Errorf("failed to download video: %w", err)
    }

    // 2. Parse filename for output directory and S3 prefix
    name := strings.TrimSuffix(filepath.Base(objectKey), filepath.Ext(objectKey))

	
    outputDir := filepath.Join(os.TempDir(), name+"_output")

    // 3. Chunk the video
    if err := ChunkVideo(localInput, outputDir); err != nil {
        _ = DeleteLocalFile(localInput)
        return fmt.Errorf("failed to chunk video: %w", err)
    }

    // 4. Upload each chunked file to S3 as {filename}/chunk_XXX.mp4
    files, err := os.ReadDir(outputDir)
    if err != nil {
        _ = DeleteLocalFile(localInput)
        _ = DeleteLocalFile(outputDir)
        return fmt.Errorf("failed to read output dir: %w", err)
    }
    for _, f := range files {
        if f.IsDir() {
            continue
        }
        chunkPath := filepath.Join(outputDir, f.Name())
        chunkKey := fmt.Sprintf("%s/%s", name, f.Name()) // {filename}/chunk_XXX.mp4
        if err := s3util.UploadFile(bucket, chunkKey, chunkPath); err != nil {
            fmt.Printf("Failed to upload chunk %s: %v\n", chunkPath, err)
        }
    }

    // 5. Clean up local files
    _ = DeleteLocalFile(localInput)
    _ = DeleteLocalFile(outputDir)

    return nil
}

// ChunkVideo splits a video into 10-second segments using ffmpeg.
func ChunkVideo(inputPath, outputDir string) error {
    segmentDuration := 10 // seconds

    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    outputPattern := filepath.Join(outputDir, "chunk_%03d.mp4")

    cmd := exec.Command(
        "ffmpeg",
        "-i", inputPath,
        "-c", "copy",
        "-map", "0",
        "-f", "segment",
        "-segment_time", fmt.Sprintf("%d", segmentDuration),
        "-reset_timestamps", "1",
        outputPattern,
    )

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("ffmpeg failed: %w", err)
    }

    return nil
}

// DeleteLocalFile deletes a file or directory at the given path.
func DeleteLocalFile(path string) error {
    return os.RemoveAll(path)
}