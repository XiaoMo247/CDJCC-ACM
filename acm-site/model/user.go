package model

import (
	"time"
)

type User struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	StudentNumber string      `gorm:"type:varchar(20);uniqueIndex;not null" json:"student_number"` // 学号，唯一
	Password      string      `gorm:"not null" json:"-"`                                           // 密码（加密后）
	Username      string      `gorm:"not null" json:"username"`                                    // 昵称/姓名
	CreatedAt     time.Time   `json:"created_at"`
	JoinApply     []JoinApply `gorm:"foreignKey:UserID"` // 一对多关联
}
