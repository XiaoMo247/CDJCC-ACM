package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

var GlobalConfig *Config

func InitConfig() {
	data, err := os.ReadFile("config/config.yml")
	if err != nil {
		panic("读取配置文件失败: " + err.Error())
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic("解析配置文件失败: " + err.Error())
	}

	// 优先使用环境变量覆盖配置文件
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.MySQL.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		// 这里简化处理，生产环境应该转换为int
		// cfg.MySQL.Port = ...
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.MySQL.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.MySQL.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.MySQL.DBName = dbname
	}

	GlobalConfig = cfg
}
