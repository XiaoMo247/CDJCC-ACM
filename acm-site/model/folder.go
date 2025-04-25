package model

import (
	"time"
)

type Folder struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;unique" json:"name"` // 文件夹名称
	CreatedAt time.Time `json:"created_at"`
}
