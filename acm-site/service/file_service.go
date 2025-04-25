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
	// 1. 先根据 folderID 查出 Folder.Name
	var folder model.Folder
	if err := database.DB.First(&folder, folderID).Error; err != nil {
		return nil, err
	}

	// 2. 构建最终的上传目录：相对于当前工作目录的上级（../）下的 uploads/<folder.Name>
	uploadDir := filepath.Join("..", "uploads", folder.Name)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	// 3. 生成一个防重名的文件名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	// 4. 拼接文件完整路径
	filePath := filepath.Join(uploadDir, filename)

	// 5. 保存文件到磁盘
	if err := saveFile(file, filePath); err != nil {
		return nil, err
	}

	// 6. 写入数据库
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
	err := database.DB.First(&file, id).Error
	if err != nil {
		return err
	}
	os.Remove(file.Path)
	return database.DB.Delete(&file).Error
}

func GetFileByID(id uint) (*model.File, error) {
	var file model.File
	err := database.DB.First(&file, id).Error
	return &file, err
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
