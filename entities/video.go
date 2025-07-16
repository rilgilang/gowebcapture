package entities

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	ID          string    `gorm:"type:varchar(36);primary_key;unique" json:"id"`
	OriginalUrl string    `gorm:"type:text" json:"original_url"`
	URL         string    `gorm:"type:text" json:"url"`
	UserID      string    `gorm:"type:varchar(36)" json:"user_id"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at," sql:"DEFAULT:CURRENT_TIMESTAMP"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at," sql:"DEFAULT:CURRENT_TIMESTAMP"`
}

type VideoQueuePayload struct {
	ID     string `json:"id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	URL    string `json:"url,omitempty"`
}
