// entity/video.go
package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Filename   string    `gorm:"not null" json:"filename"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"-"`
	Comments   []Comment `gorm:"foreignKey:VideoID" json:"comments,omitempty"`
	LikeCount  int       `gorm:"default:0" json:"like_count"`
	ViewCount  uint      `gorm:"default:0" json:"view_count"`
}

func (v *Video) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}
