package service

import (
	"acm-site/database"
	"acm-site/model"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveUploadedFile(file *multipart.FileHeader, folderID uint) (*model.File, error) {
	// Ensure folder exists.
	var folder model.Folder
	if err := database.DB.First(&folder, folderID).Error; err != nil {
		return nil, err
	}

	// Store files by folder ID (stable even if folder name/path changes).
	uploadDir := filepath.Join("..", "uploads", "folders", fmt.Sprintf("%d", folder.ID))
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	if err := saveFile(file, filePath); err != nil {
		return nil, err
	}

	newFile := model.File{
		Name:     file.Filename,
		Path:     filePath,
		FolderID: folderID,
		Size:     file.Size,
	}
	if err := database.DB.Create(&newFile).Error; err != nil {
		return nil, err
	}

	return &newFile, nil
}

func GetFilesByFolderID(folderID uint) ([]model.File, error) {
	var files []model.File
	err := database.DB.Where("folder_id = ?", folderID).Find(&files).Error
	return files, err
}

func DeleteFileByID(id uint) error {
	var file model.File
	if err := database.DB.First(&file, id).Error; err != nil {
		return err
	}
	_ = os.Remove(file.Path)
	return database.DB.Delete(&file).Error
}

func GetFileByID(id uint) (*model.File, error) {
	var file model.File
	err := database.DB.First(&file, id).Error
	return &file, err
}

// IncrementFileDownloadCount 增加文件下载计数
func IncrementFileDownloadCount(id uint) error {
	return database.DB.Model(&model.File{}).Where("id = ?", id).UpdateColumn("download_count", database.DB.Raw("download_count + 1")).Error
}

// GetTopDownloadFiles 获取下载量最高的前N个文件
func GetTopDownloadFiles(limit int) ([]model.File, error) {
	var files []model.File
	err := database.DB.Order("download_count DESC").Limit(limit).Find(&files).Error
	return files, err
}

func saveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	return err
}
