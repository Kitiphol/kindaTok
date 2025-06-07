// entity/video.go
package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    Title        string    `gorm:"not null" json:"title"`
    // filename = objectKey == 
    // should store "{video_id}_{original_filename}"
    Filename     string    `gorm:"not null" json:"filename"`
    UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"` // Foreign key to User
    User         User      `gorm:"foreignKey:UserID" json:"-"`        // Belongs to User
    Comments     []Comment `gorm:"foreignKey:VideoID" json:"comments,omitempty"` // One-to-many: Video has many Comments
    LikeCount    int       `gorm:"default:0" json:"like_count"`
    ViewCount    uint      `gorm:"default:0" json:"view_count"`
    // S3URL        string    `gorm:"not null" json:"-"`
    // ThumbnailURL string    `gorm:"not null" json:"-"`
}


func (v *Video) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}
