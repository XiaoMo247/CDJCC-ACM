package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// 获取所有文件夹
func GetAllFolders() ([]model.Folder, error) {
	var folders []model.Folder
	if err := database.DB.Order("created_at desc").Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

// 添加文件夹
func CreateFolder(name string) error {
	// 文件夹名不可重复
	var count int64
	database.DB.Model(&model.Folder{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		return errors.New("文件夹已存在")
	}
	folder := model.Folder{Name: name}
	return database.DB.Create(&folder).Error
}

// 删除文件夹
func RemoveFolder(folderID uint) error {
	// 1. 获取该文件夹的详细信息
	var folder model.Folder
	if err := database.DB.First(&folder, folderID).Error; err != nil {
		return fmt.Errorf("文件夹不存在: %w", err)
	}

	// 2. 构造文件夹的物理路径
	folderPath := filepath.Join("..", "uploads", folder.Name)
	fmt.Println(folderPath)

	// 3. 删除该文件夹及其内容
	err := os.RemoveAll(folderPath)
	if err != nil {
		return fmt.Errorf("删除文件夹失败: %w", err)
	}

	// 4. 删除数据库中的文件夹记录
	if err := database.DB.Delete(&folder).Error; err != nil {
		return fmt.Errorf("删除文件夹数据库记录失败: %w", err)
	}

	return nil
}
