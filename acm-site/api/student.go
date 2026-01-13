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

// 学生登录
func StudentLogin(c *gin.Context) {
	var loginService service.StudentLoginService
	if err := c.ShouldBindJSON(&loginService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	student, err := loginService.Login(database.DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error(), "message": err.Error()})
		return
	}

	// 使用统一的 token 生成函数
	token, err := jwt.GenerateUnifiedToken(student.ID, student.StudentID, "member")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "生成 token 失败"})
		return
	}

	// 返回统一格式：token + user 对象
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "登录成功",
		"token": token,
		"user": gin.H{
			"id":         student.ID,
			"username":   student.Username,
			"student_id": student.StudentID,
			"role":       "member",
		},
	})
}

// 获取当前学生信息
func GetStudentInfo(c *gin.Context) {
	// 从统一中间件获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	var student model.TeamMember
	if err := database.DB.First(&student, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "学生不存在"})
		return
	}

	student.Password = "-"

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": student,
	})
}

// 修改学生密码
type UpdateStudentPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	Confirm     string `json:"confirm" binding:"required"`
}

func UpdateStudentPassword(c *gin.Context) {
	var req UpdateStudentPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if req.NewPassword != req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "两次密码不一致"})
		return
	}

	// 从统一中间件获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	var student model.TeamMember
	if err := database.DB.First(&student, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "学生不存在"})
		return
	}

	if !utils.CheckPassword(req.OldPassword, student.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "原密码错误"})
		return
	}

	student.Password = utils.HashPassword(req.NewPassword)
	database.DB.Save(&student)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "密码更新成功，需重新登录"})
}

// 修改学生用户名
type UpdateStudentUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

func UpdateStudentUsername(c *gin.Context) {
	var req UpdateStudentUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 从统一中间件获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	var student model.TeamMember
	if err := database.DB.First(&student, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "学生不存在"})
		return
	}

	student.Username = req.Username
	database.DB.Save(&student)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "用户名修改成功"})
}

type UpdateStudentInfoRequest struct {
	CfName string `json:"cf_name"`
	AtName string `json:"at_name"`
	NcID   string `json:"nc_id"`
}

// 修改学生的 Codeforces, AtCoder 和 Nowcoder 用户信息
func UpdateStudentInfo(c *gin.Context) {
	var req UpdateStudentInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 从统一中间件获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	var student model.TeamMember
	if err := database.DB.First(&student, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}

	// 更新学生的 Codeforces, AtCoder 和 Nowcoder 用户信息
	student.CfName = req.CfName
	student.AtName = req.AtName
	student.NcID = req.NcID

	// 保存更新后的信息
	if err := database.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "信息更新成功",
	})
}
