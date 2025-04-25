package service

import (
	"acm-site/model"
	"acm-site/utils"
	"errors"

	"gorm.io/gorm"
)

type StudentLoginService struct {
	StudentID string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (s *StudentLoginService) Login(db *gorm.DB) (*model.TeamMember, error) {
	var student model.TeamMember
	if err := db.Where("student_id = ?", s.StudentID).First(&student).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if !utils.CheckPassword(s.Password, student.Password) {
		return nil, errors.New("密码错误")
	}

	return &student, nil
}
