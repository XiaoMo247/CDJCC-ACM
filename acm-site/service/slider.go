package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"

	"gorm.io/gorm"
)

// 获取轮播图列表
func GetSliderList() ([]model.Slider, error) {
	var sliders []model.Slider
	if err := database.DB.Find(&sliders).Error; err != nil {
		return nil, err
	}
	return sliders, nil
}

// 添加轮播图
func CreateSlider(slider *model.Slider) error {
	if slider.Title == "" || slider.Image == "" {
		return errors.New("标题和图片不能为空")
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
