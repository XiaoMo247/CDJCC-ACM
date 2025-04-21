package service

import (
	"fmt"
	"strconv"

	"acm-site/database"
	"acm-site/model"
	"acm-site/utils"
)

func RegisterUserBatch(start, end int) ([]string, error) {
	if start > end {
		return nil, fmt.Errorf("起始学号不能大于结束学号")
	}

	db := database.DB
	created := []string{}

	for num := start; num <= end; num++ {
		studentNumber := strconv.Itoa(num)
		password := studentNumber[len(studentNumber)-6:] // 后六位为默认密码

		var count int64
		db.Model(&model.User{}).Where("student_number = ?", studentNumber).Count(&count)
		if count > 0 {
			continue // 已存在，跳过
		}

		user := model.User{
			StudentNumber: studentNumber,
			Password:      utils.HashPassword(password),
			Username:      "未命名用户",
		}
		if err := db.Create(&user).Error; err != nil {
			return created, err
		}
		created = append(created, studentNumber)
	}

	return created, nil
}
