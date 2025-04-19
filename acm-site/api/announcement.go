package api

import (
	"acm-site/database"
	"acm-site/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建公告
func CreateAnnouncement(c *gin.Context) {
	var req model.Announcement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "标题和内容不能为空"})
		return
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
}

// 获取所有公告
func ListAnnouncements(c *gin.Context) {
	var announcements []model.Announcement
	if err := database.DB.Order("created_at desc").Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"announcements": announcements})
}

// 删除公告
func DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.Announcement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

// 修改公告
func UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")
	var req model.Announcement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	var ann model.Announcement
	if err := database.DB.First(&ann, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "公告不存在"})
		return
	}

	ann.Title = req.Title
	ann.Content = req.Content
	if err := database.DB.Save(&ann).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "修改失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}
