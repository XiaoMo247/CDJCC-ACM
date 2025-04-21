package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"
)

// 提交申请表
func SubmitJoinApply(apply *model.JoinApply) error {
	db := database.DB

	var count int64
	err := db.Model(&model.JoinApply{}).
		Where("user_id = ? AND apply = ?", apply.UserID, apply.Apply).
		Count(&count).Error

	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("你已申请该职位，请勿重复提交")
	}

	return db.Create(apply).Error
}

// 获取用户所有申请记录
func GetMyJoinApplies(userID uint) ([]model.JoinApply, error) {
	var applies []model.JoinApply
	err := database.DB.Where("user_id = ?", userID).Find(&applies).Error
	if err != nil {
		return nil, err
	}
	return applies, nil
}

// 管理员获取所有申请记录
func GetAllJoinApplies() ([]model.JoinApply, error) {
	var applies []model.JoinApply
	err := database.DB.Order("created_at DESC").Find(&applies).Error
	if err != nil {
		return nil, err
	}
	return applies, nil
}

// 根据 ID 获取申请记录
func GetJoinApplyByID(id uint) (*model.JoinApply, error) {
	var apply model.JoinApply
	err := database.DB.First(&apply, id).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

// 通过某一申请
func ApproveJoinApply(id uint) error {
	var apply model.JoinApply

	// 先查找该申请是否存在
	if err := database.DB.First(&apply, id).Error; err != nil {
		return err
	}

	// 修改状态为“审核通过”
	apply.State = 2
	return database.DB.Save(&apply).Error
}

// RejectJoinApply 管理员审核不通过某一申请表
func RejectJoinApply(id uint) error {
	var apply model.JoinApply
	if err := database.DB.First(&apply, id).Error; err != nil {
		return err
	}

	// 设置状态为 "不通过"（1）
	apply.State = 1
	return database.DB.Save(&apply).Error
}
