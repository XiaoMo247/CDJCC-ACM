package model

import (
	"gorm.io/gorm"
)

type FAQ struct {
	gorm.Model
	Question string `gorm:"type:varchar(255);not null" json:"question"` // 问题内容
	Answer   string `gorm:"type:text;not null" json:"answer"`           // 答案内容
}
