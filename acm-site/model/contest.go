package model

import "time"

type Contest struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"` // 比赛标题
	Text      string    `gorm:"type:text" json:"text"`                   // 简要介绍或说明
	Link      string    `gorm:"type:varchar(255)" json:"link"`           // 外部链接（例如比赛平台链接）
	Time      time.Time `gorm:"not null" json:"time"`                    // 比赛时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
