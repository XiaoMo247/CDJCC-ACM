package model

import (
	"time"
)

type Slider struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Image     string    `gorm:"type:varchar(255);not null" json:"image"`
	CreatedAt time.Time `json:"created_at"`
}
