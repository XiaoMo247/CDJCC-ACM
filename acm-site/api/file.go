package api

import (
	"acm-site/service"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	folderIDStr := c.PostForm("folder_id")
	folderID, err := strconv.Atoi(folderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的文件夹ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件获取失败"})
		return
	}

	newFile, err := service.SaveUploadedFile(file, uint(folderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "文件保存失败"})
		return
	}

	fmt.Println(newFile.Name)

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"file":    newFile, // newFile 中已有 size 字段
	})
}

func ListFiles(c *gin.Context) {
	folderIDStr := c.Query("folder_id")
	folderID, err := strconv.Atoi(folderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的文件夹ID"})
		return
	}

	files, err := service.GetFilesByFolderID(uint(folderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取文件列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

func DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的文件ID"})
		return
	}

	file, err := service.GetFileByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "文件不存在"})
		return
	}

	// 增加下载计数
	if err := service.IncrementFileDownloadCount(uint(id)); err != nil {
		// 记录错误但不影响下载流程
		fmt.Printf("更新下载计数失败: %v\n", err)
	}

	// 读取文件内容
	data, err := os.ReadFile(file.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取文件失败"})
		return
	}

	// Base64 编码
	encoded := base64.StdEncoding.EncodeToString(data)

	c.JSON(http.StatusOK, gin.H{
		"name": file.Name,
		"type": "application/octet-stream",
		"data": encoded,
	})
}

func DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的文件ID"})
		return
	}

	err = service.DeleteFileByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetFileDownloadStats 获取文件下载统计（前10）
func GetFileDownloadStats(c *gin.Context) {
	files, err := service.GetTopDownloadFiles(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取统计数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": files})
}
