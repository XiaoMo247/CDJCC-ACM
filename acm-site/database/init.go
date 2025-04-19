package database

import (
	"acm-site/config"
	"acm-site/model"
	"acm-site/utils"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	c := config.GlobalConfig.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	err = DB.AutoMigrate(&model.Admin{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	err = DB.AutoMigrate(&model.Announcement{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	var count int64
	DB.Model(&model.Admin{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		admin := model.Admin{
			Username: "admin",
			Password: utils.HashPassword("admin123"),
		}
		DB.Create(&admin)
		fmt.Println("已创建初代管理员：admin / admin123")
	}
}
