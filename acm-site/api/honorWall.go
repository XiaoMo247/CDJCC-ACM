package api

import (
	"net/http"
	"strconv"

	"acm-site/database"
	"acm-site/model"

	"github.com/gin-gonic/gin"
)

// 创建荣誉记录
func CreateHonor(c *gin.Context) {
	var honor model.Honor
	if err := c.ShouldBindJSON(&honor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := database.DB.Create(&honor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": honor,
	})
}

// 删除荣誉记录
func DeleteHonor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的ID"})
		return
	}

	if err := database.DB.Delete(&model.Honor{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

// 更新荣誉记录
func UpdateHonor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的ID"})
		return
	}

	var honor model.Honor
	if err := c.ShouldBindJSON(&honor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := database.DB.Model(&model.Honor{}).Where("id = ?", id).Updates(honor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// 获取所有荣誉记录
func GetAllHonors(c *gin.Context) {
	var honors []model.Honor
	if err := database.DB.Find(&honors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": honors,
	})
}
