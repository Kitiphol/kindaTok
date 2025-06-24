package service

import (
    "WebsocketService/internal/entity"
    "WebsocketService/internal/middleware"
    "WebsocketService/internal/database"
    "github.com/gin-gonic/gin"
    "net/http"
    "github.com/google/uuid"
	"gorm.io/gorm"
)

// Add or update a comment on a video
func AddOrUpdateComment(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    var req struct {
        Content string `json:"content" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment"})
        return
    }

    comment := entity.Comment{
        Content: req.Content,
        UserID:  userID,
        VideoID: videoID,
    }

    // Save comment to DB
    if err := database.DB.Create(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Comment added/updated", "comment": comment})
}

func AddOrUpdateLike(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoID, err := uuid.Parse(c.Param("videoID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }

    // Check existing
    var like entity.Like
    if err := database.DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err == nil {
        // Already liked, return current state
        var video entity.Video
        database.DB.First(&video, "id = ?", videoID)
        c.JSON(http.StatusOK, gin.H{"likes": video.TotalLikeCount, "hasLiked": true})
        return
    }

    // Create like
    newLike := entity.Like{UserID: userID, VideoID: videoID}
    if err := database.DB.Create(&newLike).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add like"})
        return
    }
    // Increment count
    if err := database.DB.Model(&entity.Video{}).
        Where("id = ?", videoID).
        UpdateColumn("total_like_count", gorm.Expr("total_like_count + 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment like count"})
        return
    }

    // Return updated
    var video entity.Video
    database.DB.First(&video, "id = ?", videoID)
    c.JSON(http.StatusOK, gin.H{"likes": video.TotalLikeCount, "hasLiked": true})
}

// Add or update a view on a video
func AddOrUpdateView(c *gin.Context) {
    _, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }

    // Increment view count in Video
    if err := database.DB.Model(&entity.Video{}).Where("id = ?", videoID).UpdateColumn("total_view_count", gorm.Expr("total_view_count + ?", 1)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update view count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "View added/updated"})
}

// Get total likes for a video
func GetTotalLikes(c *gin.Context) {
    videoID, err := uuid.Parse(c.Param("videoID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    // fetch video
    var video entity.Video
    if err := database.DB.First(&video, "id = ?", videoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
        return
    }
    // check user like
    userID, _ := middleware.GetUserIDFromContext(c)
    var count int64
    if userID != uuid.Nil {
        database.DB.Model(&entity.Like{}).
            Where("user_id = ? AND video_id = ?", userID, videoID).
            Count(&count)
    }
    c.JSON(http.StatusOK, gin.H{"likes": video.TotalLikeCount, "hasLiked": count > 0})
}


// Get total views for a video
func GetTotalViews(c *gin.Context) {
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    var video entity.Video
    if err := database.DB.First(&video, "id = ?", videoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"views": video.TotalViewCount})
}

// List all comments for a video
func ListComments(c *gin.Context) {
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    var comments []entity.Comment
    if err := database.DB.Where("video_id = ?", videoID).Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"comments": comments})
}


// Delete a like on a video (unlike)
// DeleteLike removes a user's like and decrements the count
func DeleteLike(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoID, err := uuid.Parse(c.Param("videoID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }

    // Delete record
    res := database.DB.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&entity.Like{})
    if res.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete like"})
        return
    }
    if res.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
        return
    }
    // Decrement count
    if err := database.DB.Model(&entity.Video{}).
        Where("id = ?", videoID).
        UpdateColumn("total_like_count", gorm.Expr("GREATEST(total_like_count - 1, 0)")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrement like count"})
        return
    }

    var video entity.Video
    database.DB.First(&video, "id = ?", videoID)
    c.JSON(http.StatusOK, gin.H{"likes": video.TotalLikeCount, "hasLiked": false})
}

// Delete a comment for a video (user specified)
func DeleteComment(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    commentIDStr := c.Param("commentID")
    commentID, err := uuid.Parse(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid commentID"})
        return
    }

    // Only allow the user who created the comment to delete it
    var comment entity.Comment
    if err := database.DB.First(&comment, "id = ? AND video_id = ? AND user_id = ?", commentID, videoID, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or not authorized"})
        return
    }
    if err := database.DB.Delete(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}