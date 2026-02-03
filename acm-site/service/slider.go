package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"

	"gorm.io/gorm"
)

// 获取轮播图列表（预加载关联的Image）
func GetSliderList() ([]model.Slider, error) {
	var sliders []model.Slider
	// Preload关联的Image对象，只查询is_active为true的轮播图，按order排序
	if err := database.DB.Where("is_active = ?", true).
		Preload("ImageObj").
		Order("`order` ASC, created_at DESC").
		Find(&sliders).Error; err != nil {
		return nil, err
	}
	return sliders, nil
}

// 添加轮播图
func CreateSlider(slider *model.Slider) error {
	// 验证：要么有ImageID，要么有Image URL
	if slider.ImageID == nil && slider.Image == "" {
		return errors.New("必须提供图片ID或图片URL")
	}
	if slider.Title == "" {
		return errors.New("标题不能为空")
	}
	if err := database.DB.Create(slider).Error; err != nil {
		return err
	}
	return nil
}

// 删除轮播图
func DeleteSlider(id uint) error {
	if err := database.DB.Delete(&model.Slider{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("轮播图不存在")
		}
		return err
	}
	return nil
}
