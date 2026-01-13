package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MySQL    MySQLConfig     `yaml:"mysql"`
	JWT      JWTConfig       `yaml:"jwt"`
	CORS     CORSConfig      `yaml:"cors"`
	Security SecurityConfig  `yaml:"security"`
}

type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	Charset         string `yaml:"charset"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
}

type JWTConfig struct {
	SecretKey       string `yaml:"secret_key"`
	ExpirationHours int    `yaml:"expiration_hours"`
}

type CORSConfig struct {
	AllowedOrigins string `yaml:"allowed_origins"`
}

type SecurityConfig struct {
	SuperAdminIDs string `yaml:"super_admin_ids"`
}

var GlobalConfig *Config

func InitConfig() {
	// 首先尝试加载配置文件（如果存在）
	var cfg *Config

	// 如果配置文件存在，从文件加载
	if data, err := os.ReadFile("config/config.yml"); err == nil {
		cfg = &Config{}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			panic("解析配置文件失败: " + err.Error())
		}
	} else {
		// 如果文件不存在，创建默认配置
		cfg = &Config{
			MySQL: MySQLConfig{
				Host:            "127.0.0.1",
				Port:            3306,
				User:            "root",
				Password:        "123456",
				DBName:          "acm_site",
				Charset:         "utf8mb4",
				MaxIdleConns:    10,
				MaxOpenConns:    100,
			},
			JWT: JWTConfig{
				SecretKey:       "default_dev_secret_key_please_change_in_production",
				ExpirationHours: 24,
			},
			CORS: CORSConfig{
				AllowedOrigins: "http://localhost:5173",
			},
			Security: SecurityConfig{
				SuperAdminIDs: "1",
			},
		}
	}

	// 环境变量覆盖配置（优先级最高）
	// 数据库配置
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.MySQL.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.MySQL.Port = port
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.MySQL.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.MySQL.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.MySQL.DBName = v
	}
	if v := os.Getenv("DB_CHARSET"); v != "" {
		cfg.MySQL.Charset = v
	}
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		if conns, err := strconv.Atoi(v); err == nil {
			cfg.MySQL.MaxIdleConns = conns
		}
	}
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		if conns, err := strconv.Atoi(v); err == nil {
			cfg.MySQL.MaxOpenConns = conns
		}
	}

	// JWT 配置
	if v := os.Getenv("JWT_SECRET_KEY"); v != "" {
		cfg.JWT.SecretKey = v
	}
	if v := os.Getenv("JWT_EXPIRATION_HOURS"); v != "" {
		if hours, err := strconv.Atoi(v); err == nil {
			cfg.JWT.ExpirationHours = hours
		}
	}

	// CORS 配置
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		cfg.CORS.AllowedOrigins = v
	}

	// 安全配置
	if v := os.Getenv("SUPER_ADMIN_IDS"); v != "" {
		cfg.Security.SuperAdminIDs = v
	}

	GlobalConfig = cfg
}

// GetMySQLConfig 获取数据库配置
func GetMySQLConfig() MySQLConfig {
	if GlobalConfig == nil {
		InitConfig()
	}
	return GlobalConfig.MySQL
}

// GetJWTConfig 获取 JWT 配置
func GetJWTConfig() JWTConfig {
	if GlobalConfig == nil {
		InitConfig()
	}
	return GlobalConfig.JWT
}

// GetCORSConfig 获取 CORS 配置
func GetCORSConfig() CORSConfig {
	if GlobalConfig == nil {
		InitConfig()
	}
	return GlobalConfig.CORS
}

// GetSecurityConfig 获取安全配置
func GetSecurityConfig() SecurityConfig {
	if GlobalConfig == nil {
		InitConfig()
	}
	return GlobalConfig.Security
}
