package api

import (
	"acm-site/model"
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadContest(c *gin.Context) {
	var contest model.Contest
	if err := c.ShouldBindJSON(&contest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数解析失败",
		})
		return
	}

	if err := service.UploadContest(&contest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "比赛信息保存失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "比赛信息上传成功",
	})
}

func DeleteContest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的比赛ID",
		})
		return
	}

	if err := service.DeleteContest(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除比赛失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

func GetAllContests(c *gin.Context) {
	contests, err := service.GetAllContests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取比赛信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    contests,
	})
}
