package service

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
	"encoding/json"

    "VideoChuncker/internal/s3util"
	"github.com/RichardKnop/machinery/v2/tasks"
	"VideoChuncker/internal/machineryUtil"
	"github.com/RichardKnop/machinery/v2"
)


func ChunkAndUpload(server *machinery.Server, bucket, objectKey string) error {

    // Download the video from S3 to a local temp file
    localInput := filepath.Join(os.TempDir(), filepath.Base(objectKey))

	
    if err := s3util.DownloadFromR2(bucket, objectKey, localInput); err != nil {
        return fmt.Errorf("failed to download video: %w", err)
    }

    // Parse filename for output directory and S3 prefix
    name := strings.TrimSuffix(filepath.Base(objectKey), filepath.Ext(objectKey))


    outputDir := filepath.Join(os.TempDir(), name+"_output")

    //Chunk the video
    if err := ChunkVideo(localInput, outputDir); err != nil {
        _ = DeleteLocalFile(localInput)
        return fmt.Errorf("failed to chunk video: %w", err)
    }

    // Upload each chunked file to S3 as {filename}/chunk_XXX.mp4
    files, err := os.ReadDir(outputDir)
    if err != nil {
        _ = DeleteLocalFile(localInput)
        _ = DeleteLocalFile(outputDir)
        return fmt.Errorf("failed to read output dir: %w", err)
    }

	// server, err := machineryutil.CreateMachineryServer()
    // if err != nil {
    //     return fmt.Errorf("failed to create machinery server: %w", err)
    // }

	var chunkKeys []string

    for _, f := range files {
        if f.IsDir() {
            continue
        }

		 // tmp/{filename}_output/chunk_XXX.mp4
        chunkPath := filepath.Join(outputDir, f.Name())

		// {filename}/chunk_XXX.mp4
        chunkKey := fmt.Sprintf("%s/%s", name, f.Name()) 
        if err := s3util.UploadFile(bucket, chunkKey, chunkPath); err != nil {
            fmt.Printf("Failed to upload chunk %s: %v\n", chunkPath, err)
        }
		chunkKeys = append(chunkKeys, chunkKey)


		// // In case for parallel processing, send each chunk to the convertor
		// err := machineryutil.SendTaskToConvertor(
		// 	server,
		// 	"ConvertChunk",
		// 	tasks.Arg{Type: "string", Value: bucket},
		// 	tasks.Arg{Type: "string", Value: chunkKey},
		// )

		// if err != nil {
		// 	fmt.Printf("Failed to send convert task for %s: %v\n", chunkKey, err)
		// }

    }


	// Serialize chunkKeys to JSON
	// since Machinery don't support slices and structs
    chunkKeysJSON, err := json.Marshal(chunkKeys)
    if err != nil {
        return fmt.Errorf("failed to marshal chunk keys: %w", err)
    }

    // Send a single task with all chunk keys (one whole batch)
    err = machineryutil.SendTaskToConvertor(
        server,
        "ConvertChunks",
        tasks.Arg{Type: "string", Value: bucket},
        tasks.Arg{Type: "string", Value: string(chunkKeysJSON)},
    )


	if err := s3util.DeleteVideoFromR2(bucket, objectKey); err != nil {
        fmt.Printf("Failed to delete original video from R2: %v\n", err)
    }

    // Clean up local files (temp dir)
    _ = DeleteLocalFile(localInput)
    _ = DeleteLocalFile(outputDir)

    return nil
}

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