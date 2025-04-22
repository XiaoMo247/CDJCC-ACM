package api

import (
	"acm-site/model"
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 提交申请表
func SubmitJoinApplyHandler(c *gin.Context) {
	// 从上下文中获取 userID（由中间件注入）
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDAny.(uint)

	var apply model.JoinApply
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数绑定失败"})
		return
	}

	// 设置 userID
	apply.UserID = userID

	if err := service.SubmitJoinApply(&apply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "申请已提交"})
}

// 获取用户所有申请记录
func GetMyJoinApplyHandler(c *gin.Context) {
	userIDAny, exists := c.Get("userID") // JWT 中间件已注入 user_id
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	userID := userIDAny.(uint)

	applies, err := service.GetMyJoinApplies(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": applies,
	})
}

// 管理员查看所有申请
func GetAllJoinAppliesHandler(c *gin.Context) {
	applies, err := service.GetAllJoinApplies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applies": applies,
	})
}

// 查看某个申请表详情（管理员）
func GetJoinApplyByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 参数无效"})
		return
	}

	apply, err := service.GetJoinApplyByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该申请记录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apply": apply,
	})
}

// 审核通过某条申请
func ApproveJoinApplyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的 ID"})
		return
	}

	err = service.ApproveJoinApply(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "审核失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "审核通过"})
}

// RejectJoinApply 管理员审核不通过某一申请表
func RejectJoinApply(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	// 调用服务层方法更新申请状态为 "不通过"
	err = service.RejectJoinApply(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "审核不通过成功"})
}
