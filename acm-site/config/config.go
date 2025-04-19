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

	GlobalConfig = cfg
}
