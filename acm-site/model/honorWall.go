package model

type Honor struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`    // 姓名
	Year    int    `json:"year"`    // 获奖年份
	Contest string `json:"contest"` // 赛事类型
	Level   string `json:"level"`   // 获奖级别
	Grade   int    `json:"grade"`   // 获奖时年级
	Major   string `json:"major"`   // 专业名称
}
