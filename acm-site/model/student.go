package model

import (
	"time"
)

type TeamMember struct {
	ID        uint   `gorm:"primaryKey"`
	StudentID string `gorm:"unique;not null"` // 学号，用于登录
	Password  string `gorm:"not null"`        // 密码（加密存储）
	Username  string `gorm:"not null"`        // 显示名称

	CfName   string // Codeforces 用户名
	CfRating string // Codeforces 当前分数

	AtName   string // AtCoder 用户名
	AtRating string // AtCoder 当前分数

	NcID     string // Nowcoder 用户 ID
	NcRating string // 牛客当前分数

	Rating string // 训练赛分数
	Count  int    // 考勤次数

	CreatedAt time.Time
	UpdatedAt time.Time
}
