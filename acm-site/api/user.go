package api

import (
	"acm-site/database"
	"acm-site/model"
	"acm-site/service"
	"acm-site/utils"
	"acm-site/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	StudentNumber string `json:"student_number"`
	Password      string `json:"password"`
}

func UserLoginHandler(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数绑定失败"})
		return
	}

	var user model.User
	err := database.DB.Where("student_number = ?", req.StudentNumber).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "用户不存在"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "密码错误"})
		return
	}

	token, err := jwt.GenerateUserToken(user.ID, user.StudentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Token 生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": token,
		"user": gin.H{
			"id":             user.ID,
			"student_number": user.StudentNumber,
			"username":       user.Username,
		},
	})
}

// 修改用户名
func ChangeUsername(c *gin.Context) {
	var req struct {
		NewUsername string `json:"new_username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.NewUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	if err := service.ChangeUsername(userID.(uint), req.NewUsername); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "修改失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户名修改成功"})
}

// 修改密码
func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	err := service.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

func GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
		return
	}

	userInfo, err := service.GetUserInfo(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功",
		"data": userInfo,
	})
}
