package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Admin    AdminConfig    `mapstructure:"admin"`
	Database DatabaseConfig `mapstructure:"database"`
	Upload   UploadConfig   `mapstructure:"upload"`
	XSS      XSSConfig      `mapstructure:"xss"`
	Comment  CommentConfig  `mapstructure:"comment"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// AdminConfig 后台配置
type AdminConfig struct {
	Path     string `mapstructure:"path"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`      // sqlite 或 postgres
	Path    string `mapstructure:"path"`      // SQLite 文件路径
	Host    string `mapstructure:"host"`     // PostgreSQL 主机
	Port    int    `mapstructure:"port"`     // PostgreSQL 端口
	User    string `mapstructure:"user"`     // PostgreSQL 用户名
	Password string `mapstructure:"password"` // PostgreSQL 密码
	DBName  string `mapstructure:"dbname"`   // PostgreSQL 数据库名
	SSLMode string `mapstructure:"sslmode"` // PostgreSQL SSL模式
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path    string `mapstructure:"path"`
	MaxSize int    `mapstructure:"max_size"`
}

// XSSConfig XSS清洗配置
type XSSConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// CommentConfig 评论/留言配置
type CommentConfig struct {
	RateLimit    int      `mapstructure:"rate_limit"`    // 频率限制（秒）
	AutoApprove  bool     `mapstructure:"auto_approve"` // 是否自动审核通过
	BlockedWords []string `mapstructure:"blocked_words"` // 敏感词列表
	AdminEmails  []string `mapstructure:"admin_emails"` // 管理员邮箱（不可注册）
}

var globalConfig *Config

// Load 加载配置文件
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// Address 获取服务器地址
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
