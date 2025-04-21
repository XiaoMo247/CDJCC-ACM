package model

import (
	"time"
)

type JoinApply struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	QQNumber      string    `gorm:"not null" json:"qq_number"`      // QQ号
	StudentNumber string    `gorm:"not null" json:"student_number"` // 学号
	Apply         int       `gorm:"not null" json:"apply"`          // 申请职位
	Name          string    `gorm:"not null" json:"name"`           // 真实姓名
	Text          string    `gorm:"type:text" json:"text"`          // 自我介绍
	State         int       `gorm:"default:0" json:"state"`         // 状态：0=待审核，1=不通过，2=通过
	UserID        uint      `gorm:"not null;index" json:"user_id"`  // 外键关联 User
	CreatedAt     time.Time `json:"created_at"`
}
