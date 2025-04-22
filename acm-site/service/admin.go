package service

import (
	"fmt"
	"strconv"

	"acm-site/database"
	"acm-site/model"
	"acm-site/utils"
)

func RegisterUserBatch(startStr, endStr string) ([]string, error) {
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return nil, fmt.Errorf("起始学号格式错误")
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		return nil, fmt.Errorf("结束学号格式错误")
	}

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

// 获取所有用户
func GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsers(keyword string, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * limit
	db := database.DB.Model(&model.User{})

	if keyword != "" {
		db = db.Where("student_number LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset(offset).Limit(limit).Order("id DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func DeleteUserByID(id string) error {
	return database.DB.Delete(&model.User{}, id).Error
}

func ResetUserPassword(id string) error {
	var user model.User
	db := database.DB

	if err := db.First(&user, id).Error; err != nil {
		return err
	}

	defaultPassword := user.StudentNumber[len(user.StudentNumber)-6:] // 默认密码为学号后6位
	hashed := utils.HashPassword(defaultPassword)

	return db.Model(&user).Update("password", hashed).Error
}
