package service

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"acm-site/database"
	"acm-site/model"
	"acm-site/utils"

	"gorm.io/gorm"
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

// 根据学号获取团队成员
func GetTeamMemberByStudentID(studentID string) (*model.TeamMember, error) {
	var member model.TeamMember
	if err := database.DB.Where("student_id = ?", studentID).First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

// 删除学生
func DeleteStudent(studentID string) error {
	var member model.TeamMember
	if err := database.DB.Where("student_id = ?", studentID).First(&member).Error; err != nil {
		return err
	}

	// 删除学生
	return database.DB.Delete(&member).Error
}

// 更新密码
func UpdatePassword(studentID, newPassword string) error {
	var member model.TeamMember
	if err := database.DB.Where("student_id = ?", studentID).First(&member).Error; err != nil {
		return err
	}

	// 使用 bcrypt 加密密码
	encryptedPassword := utils.HashPassword(newPassword)

	member.Password = encryptedPassword
	return database.DB.Save(&member).Error
}

// UpdateRatings 更新学生的评分数据
func UpdateRatings(studentID, codeforcesRating, atcoderRating, nowcoderRating string) error {
	var student model.TeamMember
	// 查找学生信息
	if err := database.DB.First(&student, "student_id = ?", studentID).Error; err != nil {
		return errors.New("学生未找到")
	}

	// 更新学生的评分数据
	student.CfRating = codeforcesRating
	student.AtRating = atcoderRating
	student.NcRating = nowcoderRating

	// 保存更新后的数据
	if err := database.DB.Save(&student).Error; err != nil {
		return errors.New("更新数据库失败")
	}

	return nil
}

func GetCodeForcesRating(member *model.TeamMember, wg *sync.WaitGroup) {
	defer wg.Done()
	member.CfRating, _ = utils.FetchCodeforcesRating(member.CfName)
}

func GetAtCoderRating(member *model.TeamMember, wg *sync.WaitGroup) {
	defer wg.Done()
	member.AtRating, _ = utils.FetchAtCoderRating(member.AtName)
}

func GetNowCoderRating(member *model.TeamMember, wg *sync.WaitGroup) {
	defer wg.Done()
	member.NcRating, _ = utils.FetchNowcoderRating(member.NcID)
}

func UpdateRatingsByMember(member *model.TeamMember) error {
	var wg sync.WaitGroup
	wg.Add(3)

	go GetCodeForcesRating(member, &wg)
	go GetAtCoderRating(member, &wg)
	go GetNowCoderRating(member, &wg)

	wg.Wait() // 等待所有评分更新完成

	return database.DB.Save(member).Error
}
