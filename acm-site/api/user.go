package api

import (
	"acm-site/database"
	"acm-site/model"
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
