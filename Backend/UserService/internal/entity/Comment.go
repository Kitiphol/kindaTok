package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Content string    `gorm:"not null" json:"content"`
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	VideoID uuid.UUID `gorm:"type:uuid;not null" json:"video_id"`
	User    User      `gorm:"foreignKey:UserID" json:"-"`
	Video   Video     `gorm:"foreignKey:VideoID" json:"-"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
