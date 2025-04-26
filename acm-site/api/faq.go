package api

import (
	"acm-site/model"
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加 FAQ
func CreateFAQ(c *gin.Context) {
	var faq model.FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数绑定失败"})
		return
	}

	if err := service.CreateFAQ(&faq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "FAQ 添加成功"})
}

// 删除 FAQ
func DeleteFAQ(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	if err := service.DeleteFAQ(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "FAQ 删除成功"})
}

// GetFAQList 获取 FAQ 列表
func GetFAQList(c *gin.Context) {
	faqs, err := service.GetFAQList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取 FAQ 列表失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    faqs,
	})
}
