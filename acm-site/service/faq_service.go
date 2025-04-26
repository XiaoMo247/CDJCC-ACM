package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"
)

// 添加 FAQ
func CreateFAQ(faq *model.FAQ) error {
	if faq.Question == "" || faq.Answer == "" {
		return errors.New("问题和答案不能为空")
	}
	return database.DB.Create(faq).Error
}

// 删除 FAQ
func DeleteFAQ(id uint) error {
	return database.DB.Delete(&model.FAQ{}, id).Error
}

// GetFAQList 获取 FAQ 列表
func GetFAQList() ([]model.FAQ, error) {
	var faqs []model.FAQ
	if err := database.DB.Order("created_at DESC").Find(&faqs).Error; err != nil {
		return nil, err
	}
	return faqs, nil
}
