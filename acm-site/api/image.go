package api

import (
	"acm-site/database"
	"acm-site/model"
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请选择图片文件"})
		return
	}

	// 获取上传者信息（从token中获取）
	username, exists := c.Get("username")
	if !exists {
		username = "unknown"
	}

	// 获取可选参数
	description := c.PostForm("description")
	isSlider := c.PostForm("is_slider") == "true"

	// 保存图片
	image, err := service.SaveImage(file, username.(string), description, isSlider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"data":    image,
	})
}

// ListImages 获取图片列表
func ListImages(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 过滤参数
	isSliderFilter := c.Query("is_slider")

	// 构建查询
	query := database.DB.Model(&model.Image{})
	if isSliderFilter != "" {
		isSlider := isSliderFilter == "true"
		query = query.Where("is_slider = ?", isSlider)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询数据
	var images []model.Image
	if err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data": gin.H{
			"images": images,
			"total":  total,
			"page":   page,
			"limit":  limit,
		},
	})
}

// UpdateImage 更新图片信息
func UpdateImage(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的图片ID"})
		return
	}

	var req struct {
		IsSlider    *bool   `json:"is_slider"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}

	if err := service.UpdateImageFlags(uint(imageID), req.IsSlider, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteImage 删除图片
func DeleteImage(c *gin.Context) {
	imageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的图片ID"})
		return
	}

	if err := service.DeleteImage(uint(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
