package api

import (
	"acm-site/model"
	"acm-site/service"
	"net/http"
	"strconv"
	"os"
	"path/filepath"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// 初始化：确保 uploads/sliders/ 目录存在。
// 注意：不要在 init() 里创建 gin.Engine（会导致重复的 Gin 启动日志，而且不会被 main.go 的路由使用）。
func init() {
	if err := os.MkdirAll("uploads/sliders/", os.ModePerm); err != nil {
		panic("无法创建 uploads/sliders/ 目录: " + err.Error())
	}
}

// 获取轮播图列表
func ListSliders(c *gin.Context) {
	sliders, err := service.GetSliderList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取轮播图列表失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    sliders,
	})
}

// 添加轮播图（支持两种方式：上传图片或选择已有图片）
func AddSlider(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	imageIDStr := c.PostForm("image_id")
	orderStr := c.DefaultPostForm("order", "0")
	isActiveStr := c.DefaultPostForm("is_active", "true")

	// 解析参数
	order, _ := strconv.Atoi(orderStr)
	isActive := isActiveStr == "true"

	var imageID *uint
	var imageURL string

	// 方式1：通过image_id选择已有图片（新方式）
	if imageIDStr != "" {
		id, err := strconv.ParseUint(imageIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "无效的图片ID"})
			return
		}
		uid := uint(id)
		imageID = &uid

		// 验证图片是否存在
		var image model.Image
		if err := service.GetImageByID(uid, &image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "图片不存在"})
			return
		}
		imageURL = image.URL

	} else {
		// 方式2：上传新图片（旧方式，保持向后兼容）
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请提供image_id或上传图片文件"})
			return
		}

		// 创建保存目录
		uploadDir := "uploads/sliders/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "创建目录失败"})
			return
		}

		// 生成唯一文件名
		fileName := uuid.New().String() + filepath.Ext(file.Filename)
		filePath := filepath.Join(uploadDir, fileName)

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "文件保存失败"})
			return
		}

		imageURL = "/uploads/sliders/" + fileName
	}

	// 构造 Slider 对象
	slider := model.Slider{
		Title:    title,
		Content:  content,
		Image:    imageURL, // 保留旧字段
		ImageID:  imageID,  // 新字段
		Order:    order,
		IsActive: isActive,
	}

	if err := service.CreateSlider(&slider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "轮播图添加成功",
		"data":    slider,
	})
}

// 删除轮播图
func DeleteSlider(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	if err := service.DeleteSlider(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "轮播图删除成功"})
}
