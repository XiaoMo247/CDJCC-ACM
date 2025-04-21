package api

import (
	"acm-site/database"
	"acm-site/model"
	"acm-site/service"
	"acm-site/utils"
	"acm-site/utils/jwt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *gin.Context) {
	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求数据不合法"})
		return
	}

	var admin model.Admin
	if err := database.DB.Where("username = ?", loginDetails.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "用户名或密码错误"})
		return
	}

	if !utils.CheckPassword(loginDetails.Password, admin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "用户名或密码错误"})
		return
	}

	// 生成 token
	token, err := jwt.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成 token 失败"})
		return
	}

	// 返回 token
	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": token,
	})
}

// 添加新管理员接口
func AddAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效请求"})
		return
	}

	// 检查用户名是否已存在
	var existingAdmin model.Admin
	if err := database.DB.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"msg": "用户名已存在"})
		return
	}

	// 创建新管理员
	newAdmin := model.Admin{
		Username: req.Username,
		Password: utils.HashPassword(req.Password),
	}

	if err := database.DB.Create(&newAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建新管理员失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "新管理员已添加"})
}

func GetAdminInfo(c *gin.Context) {
	username, _ := c.Get("adminUsername")
	id, _ := c.Get("adminID")

	c.JSON(http.StatusOK, gin.H{
		"msg":      "你是合法管理员",
		"username": username,
		"id":       id,
	})
}

// 获取所有管理员列表（任意登录管理员可用）
func ListAdmins(c *gin.Context) {
	var admins []model.Admin

	if err := database.DB.Select("id", "username", "created_at").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取管理员列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"admins": admins,
	})
}

// 删除管理员接口（仅初代管理员可用）
func DeleteAdmin(c *gin.Context) {
	// 获取当前管理员 ID
	adminIDVal, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "无权限"})
		return
	}

	adminID := fmt.Sprintf("%v", adminIDVal)
	if adminID != "1" {
		c.JSON(http.StatusForbidden, gin.H{"msg": "只有初代管理员可以删除其他管理员"})
		return
	}

	// 获取 URL 参数中的 ID（要删除的管理员 ID）
	targetID := c.Param("id")
	if targetID == "1" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "初代管理员不可删除"})
		return
	}

	// 检查是否存在该管理员
	var target model.Admin
	if err := database.DB.First(&target, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "管理员不存在"})
		return
	}

	// 删除管理员
	if err := database.DB.Delete(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "管理员删除成功"})
}

type BatchRegisterRequest struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// 批量注册
func AdminBatchRegisterHandler(c *gin.Context) {
	var req BatchRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	created, err := service.RegisterUserBatch(req.Start, req.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "批量注册成功",
		"createdList": created,
		"count":       len(created),
	})
}
