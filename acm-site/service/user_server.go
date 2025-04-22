package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 修改用户名
func ChangeUsername(userID uint, newUsername string) error {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	user.Username = newUsername
	return database.DB.Save(&user).Error
}

// 修改密码
func ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = string(hashedPassword)
	return database.DB.Save(&user).Error
}

func GetUserInfo(userID uint) (*model.User, error) {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	user.Password = "" // 不返回密码
	return &user, nil
}
