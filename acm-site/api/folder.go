package api

import (
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取文件夹列表
func ListFolders(c *gin.Context) {
	folders, err := service.GetAllFolders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取文件夹失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"folders": folders})
}

// 添加文件夹
func AddFolder(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	err := service.CreateFolder(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// 删除文件夹
func DeleteFolder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的 ID"})
		return
	}
	err = service.RemoveFolder(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
