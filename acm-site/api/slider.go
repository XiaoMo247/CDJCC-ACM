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

// 初始化静态文件路由
func init() {
	// 确保 uploads/sliders/ 目录存在
	if err := os.MkdirAll("uploads/sliders/", os.ModePerm); err != nil {
		panic("无法创建 uploads/sliders/ 目录: " + err.Error())
	}

	// 配置静态文件路由
	router := gin.Default()
	router.Static("/uploads/sliders", "./uploads/sliders")
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

// 添加轮播图
func AddSlider(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件上传失败"})
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

	// 获取其他表单字段
	title := c.PostForm("title")
	content := c.PostForm("content")

	// 构造访问URL
	imageURL := "/uploads/sliders/" + fileName

	// 构造 Slider 对象
	slider := model.Slider{
		Title:   title,
		Content: content,
		Image:   imageURL,
	}

	if err := service.CreateSlider(&slider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "轮播图添加成功",
		"data": gin.H{
			"image_url": imageURL,
		},
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
