package model

import "time"

type Image struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FileName    string    `gorm:"type:varchar(255);not null" json:"file_name"`         // 原始文件名
	StoragePath string    `gorm:"type:varchar(255);not null;unique" json:"storage_path"` // 存储路径（UUID命名）
	URL         string    `gorm:"type:varchar(255);not null" json:"url"`               // 访问URL
	Size        int64     `gorm:"not null" json:"size"`                                // 文件大小（字节）
	Width       int       `gorm:"default:0" json:"width"`                              // 图片宽度
	Height      int       `gorm:"default:0" json:"height"`                             // 图片高度
	MimeType    string    `gorm:"type:varchar(100)" json:"mime_type"`                 // MIME类型
	IsSlider    bool      `gorm:"default:false" json:"is_slider"`                      // 是否用作轮播图
	Description string    `gorm:"type:text" json:"description"`                        // 图片描述（可选）
	UploadedBy  string    `gorm:"type:varchar(100)" json:"uploaded_by"`               // 上传者（管理员用户名）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
