package model

import (
	"time"
)

type Admin struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // 加密存储
	CreatedAt time.Time `json:"created_at"`
}
