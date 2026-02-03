package model

import (
	"time"
)

type Slider struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Image     string    `gorm:"type:varchar(255)" json:"image"`             // 保留旧字段用于兼容
	ImageID   *uint     `gorm:"default:null" json:"image_id"`               // 关联图片ID，允许NULL
	Order     int       `gorm:"default:0" json:"order"`                     // 轮播顺序
	IsActive  bool      `gorm:"default:true" json:"is_active"`              // 是否启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ImageObj *Image `gorm:"foreignKey:ImageID;constraint:OnDelete:SET NULL" json:"image_obj,omitempty"` // 关联的图片对象
}
