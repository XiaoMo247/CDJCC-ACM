package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// ValidateImageFile 验证图片文件类型和大小
func ValidateImageFile(file *multipart.FileHeader) error {
	// 文件类型白名单
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return errors.New("不支持的文件类型，仅支持 jpg, png, gif, webp")
	}

	// 文件大小限制：5MB
	const maxSize = 5 * 1024 * 1024
	if file.Size > maxSize {
		return errors.New("文件大小超过限制（最大5MB）")
	}

	return nil
}

// GetImageDimensions 获取图片尺寸
func GetImageDimensions(filepath string) (int, int, error) {
	img, err := imaging.Open(filepath)
	if err != nil {
		return 0, 0, err
	}
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

// SaveImage 保存图片文件并创建数据库记录
func SaveImage(file *multipart.FileHeader, uploadedBy string, description string, isSlider bool) (*model.Image, error) {
	// 1. 验证文件
	if err := ValidateImageFile(file); err != nil {
		return nil, err
	}

	// 2. 创建上传目录
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 3. 生成唯一文件名（UUID + 扩展名）
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}
	uniqueFileName := uuid.New().String() + ext
	storagePath := filepath.Join(uploadDir, uniqueFileName)

	// 4. 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	// 5. 创建目标文件
	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()

	// 6. 复制文件内容
	if _, err := dst.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 7. 获取图片尺寸
	width, height, err := GetImageDimensions(storagePath)
	if err != nil {
		// 获取尺寸失败不影响主流程，继续
		width, height = 0, 0
	}

	// 8. 生成访问URL（Windows路径转Unix风格）
	url := "/" + strings.ReplaceAll(storagePath, "\\", "/")

	// 9. 创建数据库记录
	image := &model.Image{
		FileName:    file.Filename,
		StoragePath: storagePath,
		URL:         url,
		Size:        file.Size,
		Width:       width,
		Height:      height,
		MimeType:    file.Header.Get("Content-Type"),
		IsSlider:    isSlider,
		Description: description,
		UploadedBy:  uploadedBy,
	}

	if err := database.DB.Create(image).Error; err != nil {
		// 数据库创建失败，删除已保存的文件
		os.Remove(storagePath)
		return nil, fmt.Errorf("创建数据库记录失败: %v", err)
	}

	return image, nil
}

// DeleteImage 删除图片文件和数据库记录
func DeleteImage(imageID uint) error {
	var image model.Image
	if err := database.DB.First(&image, imageID).Error; err != nil {
		return errors.New("图片不存在")
	}

	// 1. 删除文件
	if err := os.Remove(image.StoragePath); err != nil {
		// 文件删除失败不影响数据库删除
		fmt.Printf("删除文件失败: %v\n", err)
	}

	// 2. 删除数据库记录
	if err := database.DB.Delete(&image).Error; err != nil {
		return fmt.Errorf("删除数据库记录失败: %v", err)
	}

	return nil
}

// UpdateImageFlags 更新图片标记
func UpdateImageFlags(imageID uint, isSlider *bool, description *string) error {
	var image model.Image
	if err := database.DB.First(&image, imageID).Error; err != nil {
		return errors.New("图片不存在")
	}

	updates := make(map[string]interface{})
	if isSlider != nil {
		updates["is_slider"] = *isSlider
	}
	if description != nil {
		updates["description"] = *description
	}

	if err := database.DB.Model(&image).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新失败: %v", err)
	}

	return nil
}

// GetImageByID 根据ID获取图片信息
func GetImageByID(imageID uint, image *model.Image) error {
	if err := database.DB.First(image, imageID).Error; err != nil {
		return errors.New("图片不存在")
	}
	return nil
}
