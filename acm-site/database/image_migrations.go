package database

import (
	"acm-site/model"
	"fmt"
)

type columnCheckRow struct {
	Cnt int `gorm:"column:cnt"`
}

// ensureImageSchema 处理图片管理系统的数据库迁移
func ensureImageSchema() {
	if DB == nil {
		return
	}

	// 1. 迁移现有Slider数据到Image表
	migrateSliderToImage()

	// 2. 为Announcement添加view_count字段
	addAnnouncementViewCount()
}

// migrateSliderToImage 迁移现有的Slider数据到Image表
func migrateSliderToImage() {
	// 检查sliders表是否存在image_id列
	var colCheck columnCheckRow
	if err := DB.Raw(`
		SELECT COUNT(*) AS cnt
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'sliders'
		  AND COLUMN_NAME = 'image_id'
	`).Scan(&colCheck).Error; err != nil {
		fmt.Println("检查sliders表结构失败:", err)
		return
	}

	if colCheck.Cnt > 0 {
		// 已经迁移过，跳过
		return
	}

	fmt.Println("开始迁移Slider数据到Image表...")

	// 获取所有现有的Slider记录
	var oldSliders []struct {
		ID      uint
		Title   string
		Content string
		Image   string
	}
	if err := DB.Table("sliders").Find(&oldSliders).Error; err != nil {
		fmt.Println("查询旧slider数据失败:", err)
		return
	}

	if len(oldSliders) == 0 {
		fmt.Println("没有需要迁移的slider数据")
		// 直接修改表结构
		alterSliderTable()
		return
	}

	// 为每个Slider创建对应的Image记录
	for _, slider := range oldSliders {
		// 创建Image记录
		image := model.Image{
			FileName:    slider.Image,
			StoragePath: slider.Image,
			URL:         slider.Image,
			IsSlider:    true,
			Description: slider.Title,
			UploadedBy:  "system_migration",
		}

		if err := DB.Create(&image).Error; err != nil {
			fmt.Printf("创建Image记录失败 (slider_id=%d): %v\n", slider.ID, err)
			continue
		}

		fmt.Printf("已为Slider ID=%d 创建Image ID=%d\n", slider.ID, image.ID)
	}

	// 修改sliders表结构
	alterSliderTable()

	// 迁移数据到新表结构
	if err := DB.Exec(`
		UPDATE sliders s
		INNER JOIN images i ON i.storage_path = s.image
		SET s.image_id = i.id
		WHERE s.image_id IS NULL OR s.image_id = 0
	`).Error; err != nil {
		fmt.Println("迁移slider关联失败:", err)
	} else {
		fmt.Println("Slider数据迁移完成")
	}
}

// alterSliderTable 修改sliders表结构
func alterSliderTable() {
	// 添加新字段（允许NULL）
	DB.Exec("ALTER TABLE sliders ADD COLUMN IF NOT EXISTS image_id INT UNSIGNED DEFAULT NULL")
	DB.Exec("ALTER TABLE sliders ADD COLUMN IF NOT EXISTS `order` INT DEFAULT 0")
	DB.Exec("ALTER TABLE sliders ADD COLUMN IF NOT EXISTS is_active TINYINT(1) DEFAULT 1")

	// 如果image列还存在，暂时保留（用于数据迁移）
	// 迁移完成后可以手动删除
	fmt.Println("Sliders表结构已更新")
}

// addAnnouncementViewCount 为Announcement表添加view_count字段
func addAnnouncementViewCount() {
	var colCheck columnCheckRow
	if err := DB.Raw(`
		SELECT COUNT(*) AS cnt
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'announcements'
		  AND COLUMN_NAME = 'view_count'
	`).Scan(&colCheck).Error; err != nil {
		fmt.Println("检查announcements表结构失败:", err)
		return
	}

	if colCheck.Cnt == 0 {
		if err := DB.Exec("ALTER TABLE announcements ADD COLUMN view_count INT UNSIGNED DEFAULT 0").Error; err != nil {
			fmt.Println("添加view_count字段失败:", err)
		} else {
			fmt.Println("已为announcements表添加view_count字段")
		}
	}
}
