package model

import (
	"time"
)

type File struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Path          string    `gorm:"type:varchar(255);not null" json:"path"`
	FolderID      uint      `gorm:"not null" json:"folder_id"`
	Size          int64     `gorm:"not null" json:"size"`          // 文件大小（字节）
	DownloadCount uint      `gorm:"default:0" json:"download_count"` // 下载次数
	CreatedAt     time.Time `json:"created_at"`

	Folder Folder `gorm:"foreignKey:FolderID" json:"-"` // 不在 JSON 中输出 Folder 对象
}
