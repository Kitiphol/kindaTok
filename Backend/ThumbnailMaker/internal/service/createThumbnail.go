package service

import (
    
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"ThumbnailCreator/internal/s3util"


)

func CreateThumbnail(bucket, folder, segmentFilename string) error {
    // Build full S3 object key for the .ts video segment
    objectKey := fmt.Sprintf("%s/%s", folder, segmentFilename)

    // Prepare local paths
    localSegmentPath := filepath.Join(os.TempDir(), segmentFilename)
    localThumbnailPath := filepath.Join(os.TempDir(), "thumbnail.jpg")
    thumbnailKey := fmt.Sprintf("%s/thumbnail.jpg", folder)

    // Download .ts segment from R2
    if err := s3util.DownloadFromR2(bucket, objectKey, localSegmentPath); err != nil {
        return fmt.Errorf("failed to download segment: %w", err)
    }

    // Generate thumbnail using ffmpeg
    cmd := exec.Command("ffmpeg", "-y", "-i", localSegmentPath, "-ss", "00:00:01.000", "-vframes", "1", localThumbnailPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("ffmpeg failed: %w", err)
    }

    // Upload thumbnail back to R2
    if err := s3util.UploadFile(bucket, thumbnailKey, localThumbnailPath); err != nil {
        return fmt.Errorf("upload failed: %w", err)
    }

    // Clean up local files
    _ = os.Remove(localSegmentPath)
    _ = os.Remove(localThumbnailPath)

    return nil
}
