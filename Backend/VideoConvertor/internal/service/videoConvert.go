package service

import (
    "github.com/RichardKnop/machinery/v2"
	"VideoConvertor/internal/s3util"
    

	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
   
)

// Right now, convertChunks is not good
// it concatenates all chunks into one mp4 file and then converts that mp4 file into HLS segments.
// this means it chunks the video again, which is not efficient.

func ConvertChunks(ctx context.Context, server *machinery.Server, bucket string, chunkKeysJSON string) error {
    var chunkKeys []string
    if err := json.Unmarshal([]byte(chunkKeysJSON), &chunkKeys); err != nil {
        return fmt.Errorf("failed to unmarshal chunk keys: %w", err)
    }
    if len(chunkKeys) == 0 {
        return fmt.Errorf("no chunk keys provided")
    }

    // Use folder name from the first chunk as base output folder
    parts := strings.Split(chunkKeys[0], "/")
    if len(parts) < 2 {
        return fmt.Errorf("invalid chunk key format: expected 'folder/chunk_xxx.mp4'")
    }
    folderName := parts[0] // e.g. "lecture"
    hlsDir := filepath.Join(os.TempDir(), folderName)
    os.MkdirAll(hlsDir, 0755)

    var playlistLines []string
    segmentIndex := 0
    var firstSegmentKey string // Store the first segment's R2 key for thumbnail

    for _, chunkKey := range chunkKeys {
        // Download chunk locally
        localChunk := filepath.Join(os.TempDir(), filepath.Base(chunkKey))
        if err := s3util.DownloadFromR2(bucket, chunkKey, localChunk); err != nil {
            return fmt.Errorf("failed to download chunk %s: %w", chunkKey, err)
        }

        // Convert to .ts
        segmentFilename := fmt.Sprintf("segment%d.ts", segmentIndex)
        segmentPath := filepath.Join(hlsDir, segmentFilename)

        cmd := exec.Command("ffmpeg", "-i", localChunk, "-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", segmentPath)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            return fmt.Errorf("ffmpeg segment conversion failed for %s: %w", chunkKey, err)
        }

        playlistLines = append(playlistLines, "#EXTINF:10.0,", segmentFilename)

        // Upload segment to R2
        r2Key := fmt.Sprintf("%s/%s", folderName, segmentFilename)
        if err := s3util.UploadFile(bucket, r2Key, segmentPath); err != nil {
            return fmt.Errorf("failed to upload segment %s: %w", r2Key, err)
        }

        // Store first segment key for thumbnail
        if segmentIndex == 0 {
            firstSegmentKey = r2Key
        }

        // Delete local temp files
        _ = os.Remove(localChunk)
        _ = os.Remove(segmentPath)

        // Delete original chunk from R2
        if err := s3util.DeleteVideoFromR2(bucket, chunkKey); err != nil {
            return fmt.Errorf("failed to delete original chunk %s from R2: %w", chunkKey, err)
        }

        segmentIndex++
    }

    // Create playlist
    playlistPath := filepath.Join(hlsDir, "playlist.m3u8")
    playlistFile, err := os.Create(playlistPath)
    if err != nil {
        return fmt.Errorf("failed to create playlist file: %w", err)
    }
    defer playlistFile.Close()

    playlistFile.WriteString("#EXTM3U\n")
    playlistFile.WriteString("#EXT-X-VERSION:3\n")
    playlistFile.WriteString("#EXT-X-TARGETDURATION:10\n")
    playlistFile.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")
    for _, line := range playlistLines {
        playlistFile.WriteString(line + "\n")
    }
    playlistFile.WriteString("#EXT-X-ENDLIST\n")

    // Upload playlist to R2
    r2PlaylistKey := fmt.Sprintf("%s/playlist.m3u8", folderName)
    if err := s3util.UploadFile(bucket, r2PlaylistKey, playlistPath); err != nil {
        return fmt.Errorf("failed to upload playlist: %w", err)
    }

    // Send thumbnail creation task
    if server != nil && firstSegmentKey != "" {

        // thumbnailKey := fmt.Sprintf("%s/thumbnail.jpg", folderName)



        // should send bucket/{filename}/segment_0.ts
        if err := SendTaskToThumbnailCreator(server, bucket, firstSegmentKey); err != nil {
            // Log the error but don't fail the entire operation
            fmt.Errorf("Warning: failed to send thumbnail creation task: %v", err)
        }
    }

    // Cleanup
    _ = os.Remove(playlistPath)
    _ = os.RemoveAll(hlsDir)

    return nil
}
