package database

import (
	"acm-site/config"
	"acm-site/model"
	"acm-site/utils"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// generateSecurePassword 生成安全的随机密码
func generateSecurePassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic("生成随机密码失败: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func InitDB() {
	c := config.GlobalConfig.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 自动迁移所有模型
	err = DB.AutoMigrate(
		&model.Admin{},
		&model.Announcement{},
		&model.User{},
		&model.JoinApply{},
		&model.Folder{},
		&model.File{},
		&model.TeamMember{},
		&model.Contest{},
		&model.FAQ{},
		&model.Slider{},
		&model.Honor{},
	)

	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	// 初始化默认管理员
	var count int64
	DB.Model(&model.Admin{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		// 从环境变量读取初始密码，如果没有则生成随机密码
		initialPassword := os.Getenv("INITIAL_ADMIN_PASSWORD")
		if initialPassword == "" {
			initialPassword = generateSecurePassword(16)
			fmt.Println("=========================================")
			fmt.Println("🔐 初代管理员已创建")
			fmt.Println("   用户名: admin")
			fmt.Println("   密码:", initialPassword)
			fmt.Println("   ⚠️  请立即登录并修改密码！")
			fmt.Println("=========================================")
		} else {
			fmt.Println("已使用环境变量 INITIAL_ADMIN_PASSWORD 创建管理员")
		}

		admin := model.Admin{
			Username: "admin",
			Password: utils.HashPassword(initialPassword),
		}
		DB.Create(&admin)
	}
}
